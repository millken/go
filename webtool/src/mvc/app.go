package mvc

import (
	"reflect"
	"runtime"
	"strings"
	"net/http"
	"logger"
)

/*前置action*/
type cAction struct {
	controller ControllerInterface
	action string
}

type App struct {
	actions map[string][]*cAction
	View   *View
	Router *Router
}

/*
func init() {
}
*/
func NewApp() *App {
	this := new(App)
	this.actions =  make(map[string][]*cAction)
	this.View = NewView()
	this.Router = NewRouter()
	return this
}

func (this *App)AddPreAction(c ControllerInterface, a string) {
	preaction := &cAction{c, a}
	this.actions["pre"] = append(this.actions["pre"], preaction)
}

func (this *App)Run() {
	http.HandleFunc("/favicon.ico", handlerFavicon)

	http.HandleFunc("/", this.Handler)
	root := "127.0.0.1:81"

	log.Println("Http Server Started on " + root)
	err := http.ListenAndServe(root, nil)
	if err != nil {
			log.Println(err)
	}
}

func (this *App) Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("Handler crashed with error,", err)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				log.Print(file, line)
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	for _,pre := range this.actions["pre"] {
		pre.controller.SetResponse(w)
		pre.controller.SetRequest(r)

		controller := reflect.ValueOf(pre.controller)
		method := controller.MethodByName(pre.action)

		if method.Kind() != reflect.Invalid {
		//log.Println("Method:%s", method.Name)
//f := m.Func
//f.Call([]reflect.Value{reflect.ValueOf(c)})
				method.Call([]reflect.Value{})
		} else {
				http.Error(w, "Controller has no action named " + pre.action, 404)
		}
	}
	host := strings.Split(r.Host, ":")[0]
    for _, route := range append(this.Router.routes[host], this.Router.routes["*"]...) {
        if route.Match(r.URL.Path) {
	        params := r.URL.Query()
	        log.Println(r.URL.Path)
	        for key, values := range route.extractParams(r.URL.Path) {
	        		log.Printf("%s=%s",key, values)
	                params[key] = append(params[key], values...)
	        }
	        r.URL.RawQuery = params.Encode()
log.Printf("%s", r.URL.RawQuery)
            route.controller.SetResponse(w)
            route.controller.SetRequest(r)
			controller := reflect.ValueOf(route.controller)
			method := controller.MethodByName(route.action)
			method.Call([]reflect.Value{})            
            //return
        }
    }	
}


func handlerFavicon(w http.ResponseWriter, req *http.Request) {

}