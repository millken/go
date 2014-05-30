package main

import (
    "fmt"
    "github.com/mkilling/goejdb"
    "labix.org/v2/mgo/bson"
    "os"
)

func main() {
    // Create a new database file and open it
    jb, err := goejdb.Open("addressbook", JBOWRITER | JBOCREAT | JBOTRUNC)
    if err != nil {
        os.Exit(1)
    }
    // Get or create collection 'contacts'
    coll, _ := jb.CreateColl("contacts", nil)

    // Insert one record:
    // JSON: {'name' : 'Bruce', 'phone' : '333-222-333', 'age' : 58}
    rec := map[string]interface{} {"name" : "Bruce", "phone" : "333-222-333", "age" : 58}
    bsrec, _ := bson.Marshal(rec)
    coll.SaveBson(bsrec)
    fmt.Printf("\nSaved Bruce")

    // Now execute query
    res, _ := coll.Find(`{"name" : {"$begin" : "Bru"}}`) // Name starts with 'Bru' string
    fmt.Printf("\n\nRecords found: %d\n", len(res))

    // Now print the result set records
    for _, bs := range res {
        var m map[string]interface{}
        bson.Unmarshal(bs, &m)
        fmt.Println(m)
    }

    // Close database
    jb.Close()
}
