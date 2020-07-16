package rq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type rq struct {
	method string
	uri    string
	header map[string]string
	params map[string]string
	body   io.Reader
}

func DefaultRq() *rq {
	return &rq{}
}

func (r *rq) Uri(uri string) *rq {
	r.uri = uri
	return r
}

func (r *rq) SetHeader(header map[string]string) *rq {
	r.header = header
	return r
}

func (r *rq) SetParams(params map[string]string) *rq {
	r.params = params
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
		in = bytes.NewBuffer(bodyBf)
	}
	r.body = in
	return r
}

func (r *rq) SetFrom(from map[string]string, files ...*os.File) *rq {
	//写入数据
	bodyBuf := new(bytes.Buffer)
	// 创建新的写入
	sendWriter := multipart.NewWriter(bodyBuf)
	for k, v := range from {
		sendWriter.WriteField(k, v)
	}
	var err error
	for _, file := range files {
		// 创建form 上传文件
		var fileWriter io.Writer
		if fileWriter, err = sendWriter.CreateFormFile("file", "aaaaa.jpg"); err != nil {
			panic(errors.WithStack(err))
		}

		if _, err = io.Copy(fileWriter, file); err != nil {
			panic(errors.WithStack(err))
		}
		formType := sendWriter.FormDataContentType()
		r, err := http.Post("localhost:8080", formType, bodyBuf)
		if err != nil {
			panic(errors.WithStack(err))
		}
		bt, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(bt))
	}
	return r

}

func (r *rq) Set() *rq {

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
	bf, err := r.do()
	if err != nil {
		return
	}
	json.Unmarshal(bf, &res)
	return
}
func (r *rq) StringResult() (res string, err error) {
	bf, err := r.do()
	if err != nil {
		return
	}
	res = string(bf)
	return
}

func (r *rq) BufferResult() (res []byte, err error) {
	return r.do()
}

func (r *rq) do() (buff []byte, err error) {
	url := r.uri
	ps := make([]string, 0)
	// 拼接params参数
	if r.params != nil {
		for k, v := range r.params {
			ps = append(ps, fmt.Sprintf("&%s=%s", k, v))
		}
		url += strings.Join(ps, "&")
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

	client := http.DefaultClient
	rs, err := client.Do(request)
	if err != nil {
		return
	}
	defer rs.Body.Close()
	buff, err = ioutil.ReadAll(rs.Body)
	return
}
