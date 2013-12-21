package mvc
/*
import (
    "net/http"
    "reflect"
    "log"
)
type Pattern struct {

}

type Entry struct {
    W http.ResponseWriter
    R *http.Request
    Rule Rule //parse rule
    Controllers map[string]AbstractController
}

func Route(e Entry) {
    id, action := e.Rule.Parse(e.R.RequestURI)
    c, ok := e.Controllers[id]
    if !ok {
        log.Print("no route specified, will use default controller")
        c, ok = e.Controllers["*"]
        if !ok {
            log.Print("error, can not find default controller")
            return
        }
    }
    start(e.W, e.R, c, action)
}

func start(w http.ResponseWriter, r *http.Request, c AbstractController, action string) {
    c.SetW(w)
    c.SetR(r)
    t := reflect.TypeOf(c)
    m, ok := t.MethodByName(action)
    if !ok {
        log.Printf("Can not find action %s\n", action)
        return
    }

    f := m.Func
    f.Call([]reflect.Value{reflect.ValueOf(c)})
}
*/