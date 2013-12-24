package mvc

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"net/http"
	"strings"
	"logger"
)

type View struct {
	templates map[string]*template.Template
}

func NewView() *View {
	this := new(View)
	this.templates = make(map[string]*template.Template)

	return this
}

func (this *View) LoadTemplates1(tplPath string) {
    files, err := ioutil.ReadDir(tplPath)
    if err != nil { panic(err) }
    for _, file := range files {
        this.templates[file.Name()] = template.Must(template.ParseFiles(tplPath + file.Name(), "base.html"))
    }
}

func (this *View) LoadTemplates(root string, pattern string) {
        filenames, err := filepath.Glob(filepath.Join(root, pattern))
        if err != nil {
                log.Errorf("%s", err)
        }
        if len(filenames) == 0 {
                log.Errorf("view: pattern matches no files: %#q", pattern)
        }
        for i := range filenames {
        	
        	name := strings.TrimLeft(strings.TrimPrefix(filenames[i], filepath.Clean(root)), "/")
        	log.Debugf("filenames:%s, name: %s", filenames[i], name)
            t := template.New(name)
            this.templates[name], _ = t.ParseFiles(filenames[i])
        }
        //return this.ParseFiles(root, filenames...)	
}

func (this *View) Add(name, path string) *View {
	view := template.New(name)
	template.Must(view.ParseFiles(path))
	this.templates[name] = view

	return this
}

func (this *View) Render(w http.ResponseWriter, name string, data interface{}) {
	view := this.templates[name]
	log.Debugf("name: %s", name)
	err := view.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}
}
