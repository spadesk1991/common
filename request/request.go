package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type RQ struct {
	method string
	uri    string
	header map[string]string
	params map[string]string
	body   io.Reader
}

func DefaultRQ() *RQ {
	return &RQ{}
}

func (r *RQ) Uri(uri string) *RQ {
	r.uri = uri
	return r
}

func (r *RQ) SetHeader(header map[string]string) *RQ {
	r.header = header
	return r
}

func (r *RQ) SetParams(params map[string]string) *RQ {
	r.params = params
	return r
}

func (r *RQ) SetBody(body interface{}) *RQ {
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

func (r *RQ) Post() *RQ {
	r.method = "POST"
	return r
}

func (r *RQ) Get() *RQ {
	r.method = "GET"
	return r
}

func (r *RQ) Put() *RQ {
	r.method = "PUT"
	return r
}

func (r *RQ) Delete() *RQ {
	r.method = "DELETE"
	return r
}

func (r *RQ) JsonResult(res interface{}) (err error) {
	bf, err := r.do()
	if err != nil {
		return
	}
	json.Unmarshal(bf, &res)
	return
}
func (r *RQ) StringResult() (res string, err error) {
	bf, err := r.do()
	if err != nil {
		return
	}
	res = string(bf)
	return
}
func (r *RQ) do() (buff []byte, err error) {
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
