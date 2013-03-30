package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
	"log"
	"strings"
)

var db *sql.DB
var DNSBL_SUFFIX = ".dnsbl.example.com."
var DNSBL_NEIN_SUFFIX = ".nein" + DNSBL_SUFFIX

func QueryDatabase(name string) (string, error) {
	var value string

	query, err := db.Prepare("select value from dnsbl where name=?")
	if err != nil {
		log.Fatal("Cannot create prepare statement: ", err.Error())
	}

	rows, err := query.Query(name)
	if err != nil {
		log.Printf("%v error: %v\n", name, err.Error())
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&value)
		log.Println("Result=", value)
	}
	return value, nil
}

func ParseQueryName(name string) string {
	name = strings.ToLower(name)
	if !strings.HasSuffix(name, DNSBL_NEIN_SUFFIX) {
		return ""
	}
	name = name[0 : len(name)-len(DNSBL_NEIN_SUFFIX)]
	log.Println(name)
	return name
}

func ReplyNein(w dns.ResponseWriter, request *dns.Msg) {
	ans := new(dns.Msg)
	ans.SetReply(request)
	ans.Authoritative = true
	ans.Answer = make([]dns.RR, len(request.Question))

	for i, question := range request.Question {
		querypart := ParseQueryName(question.Name)
		value, err := QueryDatabase(querypart)
		if err != nil {
			// servfail
			panic("je sais pas faire servfail encore")
		}
		rr, err := dns.NewRR("nein TXT " + value)
		if err != nil {
			log.Fatal("Cannot create RR")
		}
		ans.Answer[i] = rr
	}
	w.WriteMsg(ans)
}

func InitDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("Cannot open database: ", err.Error())
	}

	_, err = db.Exec("create table if not exists dnsbl (name text, value text)")
	if err != nil {
		log.Fatal("Cannot create table: ", err.Error())
	}
	return db
}

func main() {
	db = InitDatabase("godbl.sqlite3")
	defer db.Close()

	server := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Println("Coucou")
	dns.HandleFunc("nein.dnsbl.example.com.", ReplyNein)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}

}
