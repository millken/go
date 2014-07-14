package main

import (
	"compress/gzip"
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
	Response *http.Response
}

func NewHttpx(reqUrl string) (h *Httpx) {
	headers := make(map[string]string)
	headers["User-Agent"] = headerUserAgent
	return &Httpx{
		Url:     reqUrl,
		Headers: headers,
		Method:  "GET",
		Timeout: 5,
	}
}

func (h *Httpx) AddProxy(server string) {
	h.ProxyUrl = server
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

func (h *Httpx) AddPostString(data string) {
	values, err := url.ParseQuery(data)
	if err != nil {
		h.PostData = values
		h.Method = "POST"
	}
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

func (h *Httpx) GetResponse() *http.Response {
	return h.Response
}

func disallowRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("Redirection is not allowed.")
}

//发送请求
func (h *Httpx) Send() ([]byte, error) {
	var err error
	var req *http.Request
	if h.Url == "" {
		return nil, errors.New("URL is empty")
	}

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
		req, err = http.NewRequest(h.Method, h.Url, strings.NewReader(h.PostData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(h.Method, h.Url, nil)
	}
	if err != nil {
		return nil, err
	}
	//headers
	if len(h.Headers) > 0 {
		for k, v := range h.Headers {
			if k == "Host" {
				req.Host = v
			} else {
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
		CheckRedirect: disallowRedirect,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error() + " - [ ClientIP: " + h.ClientIP + ",TargetIP: " + h.TargetIP + " ]")
	}
	h.Response = resp
	defer resp.Body.Close()
	var body []byte
	if resp.StatusCode != 200 {
		return nil, errors.New("Status != 200")
	}
	if resp.Header.Get("Content-Encoding") == "gzip" {
		var gz *gzip.Reader
		gz, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		body, err = ioutil.ReadAll(gz)
	} else {
		body, err = ioutil.ReadAll(resp.Body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil

}

/*
func main() {
	//response, err := HttpGet("http://www.baidu.com/")
	url := "https://www.baidu.com"
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
*/
