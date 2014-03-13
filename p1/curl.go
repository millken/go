package main

import (
	"bytes"
	"errors"
	curl "github.com/andelf/go-curl"
	"net/url"
	//"os"
	"fmt"
	//"time"
)

type Curlx struct {
	Debug          bool //调试开关
	Url            string
	Headers        map[string]string
	ClientIp       string //本机外网IP，可选
	TargetIp       string
	ProxyUrl       string //代理URL
	PostData       string
	ConnectTimeout int          //连接超时时间，秒
	Result         bytes.Buffer //保存结果
	Err            error        //报错信息
}

func NewCurlx() (c *Curlx) {
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36"
	reqUrl := "http://localhost/"
	return &Curlx{
		Url:            reqUrl,
		Headers:        headers,
		ConnectTimeout: 10,
		Err:            nil,
		Debug:          false,
	}
}

//设置url
func (c *Curlx) SetUrl(url string) {
	c.Url = url
}

//设置连接超时时间
func (c *Curlx) SetConnectTimeout(second int) {
	c.ConnectTimeout = second
}

//添加header
func (c *Curlx) AddHeader(key, value string) {
	c.Headers[key] = value
}

//指定出口ip
func (c *Curlx) SetClientIp(ip string) {
	c.ClientIp = ip
}

//指定目标ip
func (c *Curlx) SetTargetIp(ip string) {
	c.TargetIp = ip
}

//打开调试模式
func (c *Curlx) SetDebug() {
	c.Debug = true
}

//设置post数据
func (c *Curlx) SetPostString(data string) {
	c.PostData = data
}

//转换header map为字符数组
func (c *Curlx) HeaderString() []string {
	headers := []string{}
	for k, v := range c.Headers {
		headers = append(headers, k+": "+v)
	}
	return headers
}

//获取返回body
func (c *Curlx) GetBody() (body string, err error) {
	body = c.Result.String()
	err = c.Err
	return
}

func (c *Curlx) Send() error {
	easy := curl.EasyInit()
	defer easy.Cleanup()
	if easy != nil {
		if c.Debug == true {
			easy.Setopt(curl.OPT_VERBOSE, true)
		}
		if c.ClientIp != "" {
			easy.Setopt(curl.OPT_INTERFACE, c.ClientIp)
		}
		if c.TargetIp != "" {
			u, err := url.Parse(c.Url)
			if err != nil {
				c.Err = errors.New(err.Error() + " target ip is " + c.TargetIp)
				return c.Err
			}
			RawQuery := ""
			if u.RawQuery != "" {
				RawQuery = "?" + u.RawQuery
			}
			c.Url = fmt.Sprintf("%s://%s%s%s", u.Scheme, c.TargetIp, u.Path, RawQuery)
			c.AddHeader("Host", u.Host)
		}
		easy.Setopt(curl.OPT_URL, c.Url)
		//easy.Setopt(curl.OPT_PORT, 7891)
		easy.Setopt(curl.OPT_CONNECTTIMEOUT, c.ConnectTimeout)
		easy.Setopt(curl.OPT_HTTPHEADER, c.HeaderString())
		easy.Setopt(curl.OPT_WRITEFUNCTION, func(buf []byte, userdata interface{}) bool { c.Result.Write(buf); return true })

		if c.PostData != "" {
			easy.Setopt(curl.OPT_POSTFIELDS, c.PostData)
		}
		c.Err = easy.Perform()

	}
	return c.Err
}

func main() {
	c := NewCurlx()
	c.Url = "http://www.baidu.com/"

	c.ClientIp = "192.168.3.203"
	//c.SetTargetIp("220.181.111.86")
	c.SetPostString("a=b&c=d")
	c.AddHeader("Referer", "localhost")
	c.AddHeader("Accept", "text/html")
	c.Send()
	//c.Result.WriteTo(os.Stdout)
	body, err := c.GetBody()
	if err != nil {
		fmt.Println("Error : ", err.Error())
	} else {

		fmt.Println(body)
	}
}
