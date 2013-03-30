//http://my.oschina.net/u/126042/blog/82577
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func InitDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("Cannot open database: ", err.Error())
	}

	return db
}

func main() {
	db := InitDatabase("foo.db")
	defer db.Close()

	_, err := db.Exec("create table if not exists userinfo (uid int, username text, departname text, created date)")
	if err != nil {
		log.Fatal("Cannot create table: ", err.Error())
	}

    //插入数据
    stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	if err != nil {
		log.Fatal("Cannot prepare insert: ", err.Error())
	}

    res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	if err != nil {
		log.Fatal("Cannot insert : ", err.Error())
	}

    id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Cannot got insert id: ", err.Error())
	}

    fmt.Println(id)
    //更新数据
    stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	if err != nil {
		log.Fatal("Cannot Prepare update: ", err.Error())
	}

    res, err = stmt.Exec("astaxieupdate", id)
	if err != nil {
		log.Fatal("Cannot update: ", err.Error())
	}

    affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Cannot row affected: ", err.Error())
	}

    fmt.Println(affect)

    //查询数据
    rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		log.Fatal("Cannot select: ", err.Error())
	}

    for rows.Next() {
        var uid int
        var username string
        var department string
        var created string
        err = rows.Scan(&uid, &username, &department, &created)
        fmt.Println(uid)
        fmt.Println(username)
        fmt.Println(department)
        fmt.Println(created)
    }

    //删除数据
    stmt, err = db.Prepare("delete from userinfo where uid=?")
	if err != nil {
		log.Fatal("Cannot prepare delete: ", err.Error())
	}

    res, err = stmt.Exec(id)
	if err != nil {
		log.Fatal("Cannot delete: ", err.Error())
	}

    affect, err = res.RowsAffected()
	if err != nil {
		log.Fatal("Cannot row affected: ", err.Error())
	}

    fmt.Println(affect)

}

