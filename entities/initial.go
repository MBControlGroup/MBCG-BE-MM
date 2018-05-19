package entities

import (
    //"database/sql"
    "github.com/go-xorm/xorm"
    _ "github.com/go-sql-driver/mysql"
)

//var mydb *sql.DB
var engine *xorm.Engine

func init() {
    //https://stackoverflow.com/questions/45040319/unsupported-scan-storing-driver-value-type-uint8-into-type-time-time
    //db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
    //if err != nil {
    //    panic(err)
    //}
    //mydb = db

    //en, err := xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
    //tx, err := mydb.Begin()
    en, err := xorm.NewEngine("mysql", "mbcsdev:mbcsdev2018@(222.200.180.59:9000)/MBDB?charset=utf8")
    checkErr(err)
    engine = en
}


func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
