package rq

import (
	"fmt"
	"log"
	"testing"
)

func Test(t *testing.T) {
	DefaultSetting().Debug() // init
	//v := url.Values{"a": {"1", "2"}, "b": {"222"}}
	//fmt.Println(v.Encode())
	res, err := DefaultRq().Get().SetRetries(1).Uri("http://localhost:8080/api/os_report_data_statistics/purchase_requisition").StringResult()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
