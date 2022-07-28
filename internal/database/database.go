package database

import(
	"database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func init(){
	
	connStr := "user=postgres password=mypass dbname=productdb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    } 
    defer db.Close()
    
}