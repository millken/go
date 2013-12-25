package mvc

import (
	"reflect"
	"runtime"
	"strings"
	"net/http"
	"logger"
	"banjo"
	"fmt"
)

/*前置action*/
type cAction struct {
	controller ControllerInterface
	action string
}

type App struct {
	actions map[string][]*cAction
	Router *Router
	Config *Config
}

type Config struct {
	HttpAddr	string
	HttpPort	int
}

/*
func init() {
}
*/
func NewApp() *App {
	this := new(App)
	this.actions =  make(map[string][]*cAction)
	this.Router = NewRouter()
	this.Config = &Config{"127.0.0.1", 81}
	return this
}

func (this *App)AddPreAction(c ControllerInterface, a string) {
	preaction := &cAction{c, a}
	this.actions["pre"] = append(this.actions["pre"], preaction)
}

func (this *App)AddTemplates(dir string) {
	err := banjo.ParseTree(dir)
	if err != nil {
			log.Println(err)
	}
}

func (this *App)ServeListen(addr string, port int) {
	this.Config.HttpAddr = addr
	this.Config.HttpPort = port
}

func (this *App)ServeFile(pattern string, filename string) {
    http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filename)
    })
}

func (this *App)ServeDir(pattern string, dir string) {
	http.Handle(pattern, http.FileServer(http.Dir(dir)))
}

func (this *App)Run() {
	addr := fmt.Sprintf("%s:%d", this.Config.HttpAddr, this.Config.HttpPort)
	log.Println("Http Server Started on " + addr)
	http.HandleFunc("/", this.Handler)
	err := http.ListenAndServe(addr, nil)
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
	host := strings.Split(r.Host, ":")[0]
    for _, staticDir := range append(this.Router.staticDirs[host], this.Router.staticDirs["*"]...) {
		//static file server
		//log.Printf("r.URL.Path:%s, dir:%s", r.URL.Path, staticDir.url)
		if strings.HasPrefix(r.URL.Path, staticDir.url) {
			var file string
			if staticDir.url == "/" {
				file = staticDir.path + r.URL.Path
			} else {
				file = staticDir.path + r.URL.Path[len(staticDir.url):]
			}
			http.ServeFile(w, r, file)
			return
		}
	}
	for _,pre := range this.actions["pre"] {
		pre.controller.SetResponse(w)
		pre.controller.SetRequest(r)
		controller := reflect.ValueOf(pre.controller)
		method := controller.MethodByName(pre.action)

		if method.Kind() != reflect.Invalid {
				method.Call([]reflect.Value{})
		} else {
				http.Error(w, "Controller has no action named " + pre.action, 404)
		}
	}
    for _, route := range append(this.Router.routes[host], this.Router.routes["*"]...) {
        if route.Match(r.URL.Path) {
	        params := r.URL.Query()
	        for key, values := range route.extractParams(r.URL.Path) {
	                params[key] = append(params[key], values...)
	        }
	        r.URL.RawQuery = params.Encode()

            route.controller.SetResponse(w)
            route.controller.SetRequest(r)
            
			controller := reflect.ValueOf(route.controller)
			method := controller.MethodByName(route.action)
			method.Call([]reflect.Value{})            
            return
        }
    }	
}