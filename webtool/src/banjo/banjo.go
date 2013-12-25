// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/* Filename: banjo.go
* Author: Bryan Matsuo <bryan dot matsuo at gmail dot com>
* Created: Thu, 9 May 2013
*/

/*
Banjo is very light abstraction of the html/template library. It provides a
minimal framework for dealing with heirachical views in handrolled MVC web
apps.
*/
package banjo

import (
        "errors"
        "fmt"
        "html/template"
        "io"
        "io/ioutil"
        "os"
        "path/filepath"
        "strings"
)

var DefaultTemplate = template.New("banjo")

func Funcs(funcs template.FuncMap) {
        DefaultTemplate.Funcs(funcs)
}

func Delims(left, right string) {
        DefaultTemplate.Delims(left, right)
}

// add a template to DefaultTemplate. if name is non-empty, the template is
// surrounded in a template definition.
func Parse(name, raw string) error {
        if name != "" {
                raw = fmt.Sprintf(`{{define %q}}%s{{end}}`, name, raw)
        }
        _, err := DefaultTemplate.Parse(raw)
        return err
}

func ParseFiles(filenames ...string) error {
        _, err := DefaultTemplate.ParseFiles(filenames...)
        return err
}

func ParseGlob(pattern string) error {
        _, err := DefaultTemplate.ParseGlob(pattern)
        return err
}

// like ParseFiles and ParseGlob template names are paths relative to the root.
func ParseTree(root string) error {
        return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
                switch {
                case err != nil:
                        return err
                case info.IsDir():
                        return nil
                case !strings.HasPrefix(path, root):
                        return fmt.Errorf("unknown root: %q (%q)", path, root)
                }
                _raw, err := ioutil.ReadFile(path)
                if err != nil {
                        return err
                }
                name := path[len(root):]
                fmt.Printf("name:%s,body:%s", name, string(_raw))
                return Parse(name, string(_raw))
        })
}

// satisfied by by *template.Template.
type TemplateEngine interface {
        ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

// the engine used by NewContext. when nil, DefaultTemplate is used.
var DefaultTemplateEngine TemplateEngine

// a template context.
type Context struct {
        Data map[string]interface{}
        Engine TemplateEngine
}

// store a value in the context for use in the template.
func (context *Context) Set(key string, value interface{}) {
        context.Data[key] = value
}

// render a template.
func (context *Context) Render(wr io.Writer, name string) error {
        engine := context.Engine
        if engine == nil {
                engine = DefaultTemplate
        }
        if engine == nil {
                return errors.New("nil engine")
        }
        return engine.ExecuteTemplate(wr, name, context.Data)
}

// create a new template context using the DefaultTemplateEngine.
func NewContext() *Context {
        return &Context{
                Data: make(map[string]interface{}, 0),
        }
}

// create a new template context with an arbitrary template engine.
func NewContextEngine(engine TemplateEngine) *Context {
        c := NewContext()
        c.Engine = engine
        return c
}
