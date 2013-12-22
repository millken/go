package mvc

import (
	"html/template"
	"net/http"
)

type View struct {
	templates map[string]*template.Template
}

func NewView() *View {
	this := new(View)
	this.templates = make(map[string]*template.Template)

	return this
}

func (this *View) Add(name, path string) *View {
	view := template.New(name)
	template.Must(view.ParseFiles(path))
	this.templates[name] = view

	return this
}

func (this *View) Remove(name string) *View {

	return this
}

func (this *View) Render(w http.ResponseWriter, name string, data interface{}) {
	view := this.templates[name]
	err := view.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}
}
