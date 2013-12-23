package mvc

import (
    "net/http"
	"mime"
	//"bytes"
	"encoding/json"
	"strings"
	"strconv"
)
type Controller struct {
    Response http.ResponseWriter
    Request *http.Request
}

type ControllerInterface interface {
    SetResponse(w http.ResponseWriter)
    SetRequest(r *http.Request)
}

func (c *Controller) SetRequest(r *http.Request) {
	c.Request = r
}

func (c *Controller) SetResponse(w http.ResponseWriter) {
	c.Response = w
}



func (c *Controller) ContentType(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		c.Response.Header().Set("Content-Type", ctype)
	} else {
		c.Response.Header().Set("Content-Type", ext)
	}
}

func (c *Controller) SetHeader(hdr string, val string) {
	c.Response.Header().Set(hdr, val)
}

func (c *Controller) AddHeader(hdr string, val string) {
	c.Response.Header().Add(hdr, val)
}

func (c *Controller) Redirect(url string, code int) {
	c.Response.Header().Set("Location", url)
	c.Response.WriteHeader(code)
}

func (c *Controller) Render(contentType string, data []byte) {
	c.SetHeader("Content-Length", strconv.Itoa(len(data)))
	c.ContentType(contentType)
	c.Response.Write(data)
}

func (c *Controller) RenderHtml(content string) {
	c.Render("html", []byte(content))
}

func (c *Controller) RenderText(content string) {
	c.Render("txt", []byte(content))
}

func (c *Controller) RenderJson(data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		http.Error(c.Response, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Render("json", content)
}
