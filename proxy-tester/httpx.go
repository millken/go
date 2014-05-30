package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var headerUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36"

type Httpx struct {
	Url      string
	Headers  map[string]string
	Cookies  []*http.Cookie
	ClientIP string //本机外网IP，可选
	TargetIP string
	Method   string
	ProxyUrl string //代理URL
	PostData url.Values
	Timeout  int //超时时间，秒
}

func NewHttpx(reqUrl string) (h *Httpx) {
	headers := make(map[string]string)
	headers["User-Agent"] = headerUserAgent
	return &Httpx{
		Url:     reqUrl,
		Headers: headers,
		Method:  "GET",
		Timeout: 30,
	}
}


//添加header
func (h *Httpx) AddHeader(key, value string) {
	h.Headers[key] = value
}

//添加cookie
func (h *Httpx) AddCookie(c *http.Cookie) {
	h.Cookies = append(h.Cookies, c)
}

func (h *Httpx) SetClientIP(ip string) {
	h.ClientIP = ip
}

func (h *Httpx) SetTargetIP(ip string) {
	h.TargetIP = ip
}

func (h *Httpx) SetConnectTimeout(second int) {
	h.Timeout = second
}

//添加POST值
func (h *Httpx) AddPostValue(key string, values []string) {
	if h.PostData == nil {
		h.PostData = make(url.Values)
	}
	if values != nil {
		for _, v := range values {
			h.PostData.Add(key, v)
		}
		h.Method = "POST"
	}
}

//发送请求
func (h *Httpx) Send() (response *http.Response, err error) {
	if h.Url == "" {
		return nil, errors.New("URL is empty")
	}

	defer func() {
		if err != nil && h.ClientIP != "" {
			err = errors.New(err.Error() + " client ip is " + h.ClientIP)
		}
	}()

	var req *http.Request

	if h.TargetIP != "" {
		u, err := url.Parse(h.Url)
		if err != nil {
			err = errors.New(err.Error() + " target ip is " + h.TargetIP)
		}
		RawQuery := ""
		if u.RawQuery != "" {
			RawQuery = "?" + u.RawQuery
		}
		h.Url = fmt.Sprintf("%s://%s%s%s", u.Scheme, h.TargetIP, u.Path, RawQuery)
		h.AddHeader("Host", u.Host)
	}
	if h.Method == "POST" {
		req, _ = http.NewRequest("POST", h.Url, strings.NewReader(h.PostData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, _ = http.NewRequest(h.Method, h.Url, nil)
	}

	//headers
	if len(h.Headers) > 0 {
		for k, v := range h.Headers {
			if k == "Host" {
				req.Host = v
			}else {
			req.Header.Set(k, v)
		}
		}
	}

	//cookies
	if len(h.Cookies) > 0 {
		for _, v := range h.Cookies {
			req.AddCookie(v)
		}
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//是否使用代理
	if h.ProxyUrl != "" {
		proxy, err := url.Parse(h.ProxyUrl)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	//https://github.com/franela/goreq/blob/master/goreq.go
	//设置超时时间
	dialer := net.Dialer{
		Timeout: time.Duration(h.Timeout) * time.Second,
	}

	//是否使用指定的IP发送请求
	if h.ClientIP != "" {
		transport.Dial = func(network, address string) (net.Conn, error) {
			//本地地址  本地外网IP
			lAddr, err := net.ResolveTCPAddr(network, h.ClientIP+":0")
			if err != nil {
				return nil, err
			}
			dialer.LocalAddr = lAddr
			return dialer.Dial(network, address)
		}
	} else {
		transport.Dial = func(network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
	}

	client := &http.Client{
		Transport: transport,
	}
	response, err = client.Do(req)
	//defer response.Body.Close()
	return response, err

}

// GET请求
func HttpGet(reqUrl string) (*http.Response, error) {
	hx := NewHttpx(reqUrl)
	return hx.Send()
}

//利用指定的IP发送请求
func HttpGetFromIP(reqUrl, ipaddr string) (*http.Response, error) {
	hx := NewHttpx(reqUrl)
	hx.ClientIP = ipaddr
	return hx.Send()
}

// http GET 代理
func HttpGetFromProxy(reqUrl, proxyURL string) (*http.Response, error) {
	hx := NewHttpx(reqUrl)
	hx.ProxyUrl = proxyURL
	return hx.Send()
}

//POST请求
func HttpPost(reqUrl string, postValues map[string][]string) (*http.Response, error) {
	hx := NewHttpx(reqUrl)
	hx.Method = "POST"
	if postValues != nil {
		for k, v := range postValues {
			hx.AddPostValue(k, v)
		}
	}
	return hx.Send()
}

func main() {
	//response, err := HttpGet("http://www.baidu.com/")
	url := "http://www.baidu.com"
	httpx := NewHttpx(url)
	httpx.SetTargetIP("115.239.210.27")
	httpx.SetClientIP("192.168.3.203")
	httpx.SetConnectTimeout(1)
	httpx.AddHeader("Referer", "localhost")
	response, err := httpx.Send()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("body ", err.Error())
		}
		bodyString := string(body)
		fmt.Println(bodyString)
	}
}
