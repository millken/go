package main
//https://gist.github.com/tristanwietsma/8444cf3cb5a1ac496203

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request)

func GetOnly(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

func PostOnly(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world\n")
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.PostForm)
	io.WriteString(w, "post\n")
}

type Result struct {
	FirstName string `json:"first"`
	LastName  string `json:"last"`
}

func HandleJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(Result{"tee", "dub"})
	io.WriteString(w, string(result))
}

func main() {

	// public views
	http.HandleFunc("/", HandleIndex)

	// private views
	http.HandleFunc("/post", PostOnly(HandlePost))
	http.HandleFunc("/json", GetOnly(HandleJSON))

	log.Fatal(http.ListenAndServe(":86", nil))
}
