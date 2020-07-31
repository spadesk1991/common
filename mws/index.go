package mws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
	"time"
)

type core struct {
	secretKey    string
	accessKey    string
	sellerID     string
	mwsAuthToken string
}

type options struct {
	core
	region
}

type Option func(*options)

func newOptions(opts ...Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func newMWS(opts ...Option) *options {
	return newOptions(opts...)
}

type mwsClient struct {
	endpoint string     // Use the endpoint for your marketplace
	params   url.Values // 请求参数
	method   string     // http method
	//*options
	accessKey    string // 您的亚马逊MWS 账户由您的访问密钥编码来标识，亚马逊MWS 使用该编码来查找您的访问密钥
	secretKey    string
	sellerID     string // 卖家编号
	mwsAuthToken string // 代表亚马逊卖家对网络应用程序的特定开发商的授权

	action           string // 要对端点执行的操作，如 GetFeedSubmissionResult 操作。
	signature        string // 	验证过程中的一步，用于标识和验证请求的发送者。
	signatureMethod  string // 要用于计算签名的 HMAC 哈希算法。HmacSHA256 和 HmacSHA1 都支持哈希算法，但是亚马逊建议使用 HmacSHA256。
	signatureVersion string // 使用的签名版本。这是亚马逊MWS 特定的信息，它告诉亚马逊MWS 您使用哪种算法来生成构成签名基础的字符串。
	timestamp        string // 每个请求都必须包含请求的时间戳。根据您使用的 API 操作，您可以提供该请求的到期日期和时间，而不是时间戳。格式为 ISO 8601
	version          string // 要调用的 API 部分的版本
	userAgent        string
}

func NewMWSClient() mwsClient {
	return mwsClient{
		endpoint: "",
		params:   nil,
		method:   "",
	}
}

func (c *mwsClient) Post() *mwsClient {
	c.method = "POST"
	return c
}

func (c *mwsClient) Get() *mwsClient {
	c.method = "GET"
	return c
}

func (c *mwsClient) Put() *mwsClient {
	c.method = "PUT"
	return c
}

func (c *mwsClient) Delete() *mwsClient {
	c.method = "DELETE"
	return c
}

func (c *mwsClient) Endpoint(endpoint string) *mwsClient {
	c.endpoint = endpoint
	return c
}

func (c *mwsClient) SetParams(params url.Values) *mwsClient {
	c.params = params
	return c
}

func (c *mwsClient) do() {
	// 设置参数
	c.setParams()
	//req := http.NewRequest(c.method)
}

func (c *mwsClient) setParams() {
	c.params.Add("AWSAccessKeyId", c.accessKey)
	c.params.Add("Action", c.action)
	c.params.Add("MWSAuthToken", c.mwsAuthToken)
	c.params.Add("SellerId", c.sellerID)
	c.params.Add("SignatureMethod", c.signatureMethod)
	c.params.Add("SignatureVersion", c.signatureVersion)
	c.params.Add("SubmittedFromDate", time.Now().UTC().Format(time.RFC3339))
	c.params.Add("Timestamp", time.Now().UTC().Format(time.RFC3339))
	c.params.Add("Version", c.version)

	signature := c.sign()
	c.params.Add("Signature", signature)
}

func (c *mwsClient) calculateStringToSignV2() (str string) {
	paramsKeys := make(sort.StringSlice, 0) // 存储params key
	for k, _ := range c.params {
		paramsKeys = append(paramsKeys, k)
	}
	sort.Sort(paramsKeys)
	r := bytes.Buffer{}
	r.WriteString(c.method + "\n")
	url, _ := url.Parse(c.endpoint + "\n/\n")
	r.WriteString(url.Path)
	for i, k := range paramsKeys {
		v := c.params.Get(k)
		r.WriteString(k + "=" + v)
		if i < len(paramsKeys)-1 {
			r.WriteString("&")
		}
	}
	str = strings.ReplaceAll(r.String(), "+", "%20")
	str = strings.ReplaceAll(r.String(), "*", "%2A")
	str = strings.ReplaceAll(r.String(), "$7E", "~")
	return
}

func (c *mwsClient) sign() (signature string) {
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	str := c.calculateStringToSignV2()
	mac.Write([]byte(str))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return
}
