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
        	replace := strings.NewReplacer("\\", "/")
        	name := strings.TrimLeft(replace.Replace(strings.TrimPrefix(filenames[i], filepath.Clean(root))), "/")
        	log.Debugf("filenames:%s, name: %s", filenames[i], name)
            this.Add(name, filenames[i])
        }
        //return this.ParseFiles(root, filenames...)	
}

func (this *View) Add(name string, path string) *View {
	view := template.New(name)
	/*
	tmpl, err := view.ParseFiles(path)
	if err != nil {
		log.Errorf("Add error: %s", err)
	}
	*/
	
	this.templates[name] = template.Must(view.ParseFiles(path))

	return this
}

func (this *View) Render(w http.ResponseWriter, name string, data interface{}) {
	view := this.templates[name]
	log.Debugf("name: %s,%d", name, len(this.templates))
	err := view.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}
}
