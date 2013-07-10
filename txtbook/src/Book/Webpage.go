package Book

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05.000000000"
var db *sql.DB

func InitDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("Cannot open database: ", err.Error())
	}

	return db
}

func init() {
	db := InitDatabase("book.db")
	defer db.Close()

	_, err := db.Exec("create table if not exists book_url (url text, match text,replace text, created datetime)")
	if err != nil {
		log.Fatalf("Cannot create table: ", err.Error())
	}
}

func AddUrl(url string, match string, replace string) {
	db := InitDatabase("book.db")
	defer db.Close()
    //插入数据
    stmt, err := db.Prepare("INSERT INTO book_url(url, match, replace, created) values(?,?,?,?)")
	if err != nil {
		log.Fatalf("Cannot prepare insert: ", err.Error())
	}
	
    _, err = stmt.Exec(url, match, replace, time.Now().Format(TimeFormat[:19]))
	if err != nil {
		log.Fatalf("Cannot insert : ", err.Error())
	}
}