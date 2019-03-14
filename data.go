package main

import (
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    Uuid string
    Fname string
    Lname string
    Username string
    Email string
    Password string
}

func userExists(u *User) bool {
    db := dbConn()
    defer db.Close()
    var ps, us string
    q, err := db.Query("SELECT uname, password from user where uname=?", u.Username)
    checkErr(err)
    for q.Next() {
        q.Scan(&us, &ps)
    }
    if us == u.Username && ps == u.Password {
        return true
    }
    return false
}