package main

import (
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    delete, err := db.Query("DELETE FROM recipes WHERE uid=?", nId)
    checkErr(err)
    defer delete.Close()
    defer db.Close()
    http.Redirect(w, r, "/recipes", 302)
}