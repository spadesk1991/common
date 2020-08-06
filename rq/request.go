package rq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type rq struct {
	method     string
	uri        string
	header     map[string]string
	query      url.Values
	form       url.Values
	body       io.Reader
	bodyBf     []byte // 打印日志
	retryTime  time.Duration
	retryIndex float64
	retries    float64
}

type setting struct {
	requestId int64
	logger    logger
	debug     bool
	retryTime time.Duration
	retries   float64
	timeout   time.Duration
}

type logger interface {
	Infof(format string, args ...interface{})
}

func newSetting(log *logrus.Logger) *setting {
	return &setting{logger: log, retryTime: time.Second, retries: 4}
}

var mysetting = newSetting(logrus.New())

func DefaultSetting() *setting {
	return mysetting
}

func (s *setting) Debug() *setting {
	s.debug = true
	return s
}

func (s *setting) RetryTime(t time.Duration) *setting {
	s.retryTime = t
	return s
}

func (s *setting) Retries(count float64) *setting {
	s.retries = count
	return s
}

func (s *setting) SetLogOut(log logger) *setting {
	s.logger = log
	return s
}

func (s *setting) SetTimeOut(d time.Duration) *setting {
	s.timeout = d
	return s
}

func DefaultRq() *rq {
	return &rq{
		retryTime: mysetting.retryTime,
		retries:   mysetting.retries,
	}
}

func (r *rq) Uri(uri string) *rq {
	r.uri = uri
	return r
}

func (r *rq) SetRetryTime(t time.Duration) *rq {
	r.retryTime = t
	return r
}

func (r *rq) SetRetries(count float64) *rq {
	r.retries = count
	return r
}

func (r *rq) SetHeader(header map[string]string) *rq {
	if r.header == nil {
		r.header = map[string]string{}
	}
	for k, v := range header {
		r.header[k] = v
	}
	return r
}

func (r *rq) SetQuery(query url.Values) *rq {
	r.query = query
	return r
}

func (r *rq) SetForm(form url.Values) *rq {
	r.form = form
	return r
}

func (r *rq) SetBody(body interface{}) *rq {
	var in io.Reader
	if body != nil {
		var bodyBf []byte
		bodyBf, err := json.Marshal(body)
		if err != nil {
			panic(errors.WithStack(err))
		}
		r.bodyBf = bodyBf
		in = bytes.NewBuffer(bodyBf)
	}
	r.body = in
	return r
}

func (r *rq) SetFormFile(form map[string]string, file *os.File) *rq {
	//写入数据
	bodyBuf := new(bytes.Buffer)
	// 创建新的写入
	sendWriter := multipart.NewWriter(bodyBuf)
	for k, v := range form {
		sendWriter.WriteField(k, v)
	}
	var err error
	// 创建form 上传文件
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(rd.Intn(100))
	var fileWriter io.Writer
	if fileWriter, err = sendWriter.CreateFormFile("file", fmt.Sprintf("%s/%d%d", os.TempDir(), time.Now().UnixNano(), rd.Intn(100))); err != nil {
		panic(errors.WithStack(err))
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		panic(errors.WithStack(err))
	}
	formType := sendWriter.FormDataContentType()
	r.bodyBf = bodyBuf.Bytes()
	r.body = bodyBuf
	r.header["Content-Type"] = formType // 设置头
	return r
}

func (r *rq) Post() *rq {
	r.method = "POST"
	return r
}

func (r *rq) Get() *rq {
	r.method = "GET"
	return r
}

func (r *rq) Put() *rq {
	r.method = "PUT"
	return r
}

func (r *rq) Delete() *rq {
	r.method = "DELETE"
	return r
}

func (r *rq) JsonResult(res interface{}) (err error) {
	bf, err := r.run()
	if err != nil {
		return
	}
	json.Unmarshal(bf, &res)
	return
}

func (r *rq) StringResult() (res string, err error) {
	bf, err := r.run()
	if err != nil {
		return
	}
	res = string(bf)
	return
}

func (r *rq) BufferResult() (res []byte, err error) {
	return r.run()
}

func (r *rq) run() (buff []byte, err error) {
	var statusCode int
label:
	buff, statusCode, err = r.do()
	if statusCode >= 500 && statusCode <= 504 && r.retryIndex < r.retries {
		r.retryIndex++
		sleepTime := time.Duration(math.Pow(2, r.retryIndex)) * r.retryTime
		if mysetting.debug {
			mysetting.logger.Infof("[HTTP-RETRY] after %s", sleepTime)
		}
		time.Sleep(sleepTime)
		goto label // 重新调用
	} else { // 重置为0
		r.retryIndex = 0
	}
	return
}

func (r *rq) do() (buff []byte, statusCode int, err error) {
	statusCode = http.StatusOK
	url := r.uri
	// 解析query参数
	queryStr := r.query.Encode()

	if queryStr != "" {
		url += "?" + queryStr
	}
	if r.form != nil {
		r.header["Content-Type"] = "application/x-www-form-urlencoded"
		r.body = strings.NewReader(r.form.Encode())
	}

	request, err := http.NewRequest(r.method, url, r.body)
	if err != nil {
		return
	}
	// 设置header
	if r.header != nil {
		for k, v := range r.header {
			request.Header.Set(k, v)
		}
	}
	mysetting.requestId++
	if mysetting.debug {
		mysetting.logger.Infof("[HTTP-REQUEST] [%d] | %s | %s | %s\n", mysetting.requestId, r.method, r.uri, string(r.bodyBf))
	}
	client := http.DefaultClient
	client.Timeout = mysetting.timeout // 超时时间
	rs, err := client.Do(request)
	if err != nil {
		return
	}
	statusCode = rs.StatusCode // 返回状态码

	defer rs.Body.Close()
	buff, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		return
	}
	if rs.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("调用接口失败，[%s] | %d | %s | %s", r.method, rs.StatusCode, r.uri, string(buff)))
		return
	}
	if mysetting.debug {
		mysetting.logger.Infof("[HTTP-RESPONSE] [%d] | %s | %s | %s | %s \n", mysetting.requestId, r.method, r.uri, string(r.bodyBf), string(buff))
	}
	return
}
