package mvc

import (
    "net/http"
)
type Controller struct {
    w http.ResponseWriter
    r *http.Request
}

type ControllerInterface interface {
    SetW(w http.ResponseWriter)
    SetR(r *http.Request)
    GetW() http.ResponseWriter
    GetR() *http.Request
}


func (c *Controller) SetW(w http.ResponseWriter){
    c.w = w
}

func (c *Controller) SetR(r *http.Request){
    c.r = r
}


func (c *Controller) GetW() http.ResponseWriter{
    return c.w
}

func (c *Controller) GetR() *http.Request{
    return c.r
}