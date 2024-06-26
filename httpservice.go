package GoUtils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type HttpService struct {
	ReqObj    *http.Request
	RespObj   *http.Response
	Body      string
	CostTime  int64
	Method    string
	Headers   map[string]string
	Text      string // 响应内容字符串
	Content   []byte // 响应内容字节
	IsEchoReq bool   // 是否打印请求信息
	isDebug   bool   // 是否打印请求和结果信息

	Error error
}

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
}

// 每个http.Transport内都会维护一个自己的空闲连接池，如果每个client都创建一个新的http.Transport，就会导致底层的TCP连接无法复用。
// 如果网络请求过大，上面这种情况会导致协程数量变得非常多，导致服务不稳定。
// 为了解决这个问题，我们可以将http.Transport对象设置为全局变量，这样就可以复用连接池了。
var tr = &http.Transport{
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, //忽略https证书
	MaxIdleConnsPerHost: 2000,                                  //
	TLSHandshakeTimeout: 0 * time.Second,                       // 表示TLS 握手超时时间。这里推荐传入一个非零值0, 表示无限制
	IdleConnTimeout:     600 * time.Second,                     // 表示一个连接在空闲多久之后关闭。。这里推荐传入一个非零值0, 表示无限制
	MaxIdleConns:        0,                                     // 表示客户端对与所有host的最大空闲连接数总和。这里推荐传入一个非零值0, 表示无限制
}

type KwArgs func(hs *HttpService)

func WithHeaders(headers map[string]string) KwArgs {
	return func(hs *HttpService) {
		hs.Headers = headers
	}
}
func WithDebug(debug bool) KwArgs {
	return func(hs *HttpService) {
		hs.isDebug = debug
	}
}

func WithParams(params map[string]string) KwArgs {
	urlParams := url.Values{}
	for k, v := range params {
		urlParams.Add(k, v)
	}
	return func(hs *HttpService) {
		hs.ReqObj.URL.RawQuery = urlParams.Encode()
	}
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func (hs *HttpService) Get(url, data string, kwargs ...KwArgs) *http.Response {
	var newUrl string
	if len(data) > 0 {
		newUrl = fmt.Sprintf("%s?%s", url, data)
	} else {
		newUrl = url
	}
	text, err := hs.DoHttpRequest("GET", newUrl, "", kwargs...)
	hs.Error = err
	hs.Text = text
	if err != nil {
		log.Printf("[请求出错] %s\n", err)
	}
	return hs.RespObj
}

func (hs *HttpService) Post(url, json string, kwargs ...KwArgs) *http.Response {
	text, err := hs.DoHttpRequest("POST", url, json, kwargs...)
	hs.Text = text
	if err != nil {
		log.Printf("[请求出错] %s\n", err)
	}
	return hs.RespObj
}

func (hs *HttpService) IsPrintReq(isEchoReq bool) *HttpService {
	hs.IsEchoReq = isEchoReq
	return hs
}

func (hs *HttpService) DoHttpRequest(method, url, body string, kwargs ...KwArgs) (string, error) {
	client := hs.BuildClient()
	req, err := hs.BuildRequest(method, url, body) // 1. 构造请求对象,包含method url body信息
	if err != nil {
		return "", err
	}
	hs.ReqObj = req

	// 调用钩子函数，并对默认值进行修改,放到reqObj后，避免空指针
	for _, kwarg := range kwargs {
		kwarg(hs)
	}

	if len(hs.Headers) > 0 {
		hs.BuildRequestHeaders(hs.ReqObj, hs.Headers) // 2. 构造请求头
	} else {
		hs.BuildRequestHeaders(hs.ReqObj, defaultHeaders) // 2. 构造请求头
	}
	startTime := time.Now()              // 请求开始时间
	respObj, err := client.Do(hs.ReqObj) // 3. 发出请求
	hs.RespObj = respObj
	elapsed := time.Since(startTime).Nanoseconds() / int64(time.Millisecond) // 毫秒
	hs.CostTime = elapsed                                                    // 请求耗时绑定在实例属性上
	if hs.IsEchoReq || hs.isDebug {
		hs.PrintReqInfo(hs.ReqObj) // 4. 打印请求信息
	}
	if err != nil {
		log.Printf("[请求出错] %s\n", err)
		hs.PrintReqInfo(hs.ReqObj) // 打印请求信息
		return "", err
	}
	defer respObj.Body.Close()
	content, err := io.ReadAll(respObj.Body) // 5. 返回二进制的内容
	hs.Content = content
	if err != nil {
		log.Printf("[获取响应内容出错]%s\n", err)
		hs.PrintReqInfo(hs.ReqObj)         // 打印请求信息
		hs.PrintRespInfo(content, elapsed) // 打印结果信息
		return "", err
	}
	if hs.isDebug {
		hs.PrintRespInfo(content, elapsed) // 打印结果信息
	}
	return string(content), nil
}

func (hs *HttpService) BuildRequest(method, url string, body string) (req *http.Request, err error) {
	b := hs.BuildBody(body) // 字符串转为Reader对象
	req, err = http.NewRequest(method, url, b)
	if err != nil {
		log.Printf("BuildRequest Err %v", err)
		return nil, err
	}
	return req, nil
}

func (hs *HttpService) BuildClient() *http.Client {
	client := &http.Client{Timeout: 3 * 60 * time.Second, Transport: tr} //设置超时时间
	return client
}

func (hs *HttpService) BuildRequestHeaders(req *http.Request, headers map[string]string) *HttpService {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return hs
}

func (hs *HttpService) PrintReqInfo(req *http.Request) {
	s := fmt.Sprintf("\n    [请求Headers]：%v", req.Header) + fmt.Sprintf("\n    [请求Method]：%s", req.Method) +
		fmt.Sprintf("\n    [请求Url]：%s", req.URL) + fmt.Sprintf("\n    [请求Body]：%s", hs.Body)
	log.Println(s)
}

func (hs *HttpService) BuildBody(body string) *strings.Reader {
	hs.Body = body
	return strings.NewReader(body)
}

func (hs *HttpService) PrintRespInfo(resInfo []byte, costTime int64) *HttpService {
	costFloat := float64(costTime) / 1.0e9
	formatCostTime := fmt.Sprintf("%.3f", costFloat)
	hs.CostTime = costTime / 1e6
	r, _ := PrettyPrint(resInfo)
	s := fmt.Sprintf("\n    [响应HttpCode]：%d", hs.RespObj.StatusCode) + fmt.Sprintf("\n    [响应耗时]：%s秒",
		formatCostTime) + fmt.Sprintf("\n    [响应Body]：%s", r)
	fmt.Println(s)
	return hs
}

func (hs *HttpService) PrettyPrint(resInfo []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, resInfo, "", " "); err != nil {
		return string(resInfo), err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func (hs *HttpService) Map2String(body map[string]interface{}) string {
	return Map2JsonString(body)
}

// Json 使用方法参考https://github.com/tidwall/gjson
func (hs *HttpService) Json() gjson.Result {
	return gjson.Parse(hs.Text)
}

func (hs *HttpService) Success() bool {
	return hs.RespObj.StatusCode == 200
}

func Map2UrlValues(m map[string]interface{}) string {
	values := url.Values{}
	for k, v := range m {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return values.Encode()
}

// DoRequest 简洁版本
func DoRequest(method string, url string, requestBody []byte, requestHeaders map[string]string) (int, http.Header, string, error) {
	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, nil, "", fmt.Errorf("创建请求失败: %v", err)
	}
	for key, value := range requestHeaders {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, "", fmt.Errorf("读取响应结果失败: %v", err)
	}
	return resp.StatusCode, resp.Header, responseBody.String(), nil
}
