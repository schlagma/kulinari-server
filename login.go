package main

import (
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    redirectValue := queryValues.Get("redirect")
    name := r.FormValue("uname")
    pass := r.FormValue("password")
    u := &User{Username: name, Password: pass}
    redirect := "/?redirect=" + redirectValue
    
    if name != "" && pass != "" {
        if userExists(u) && redirectValue == "" {
            setSession(u, w)
            redirect = "/recipes"
        } else if userExists(u) && redirectValue != "" {
            setSession(u, w)
            redirect = redirectValue
        } else if u.Username != "" && redirectValue != "" {
            redirect = "/?user=" + u.Username + "&?redirect=" + redirectValue
           // setMsg(w, "message", []byte("Benutzername und/oder Passwort falsch."))
        } else if u.Username != "" && redirectValue == "" {
            redirect = "/?user=" + u.Username
           // setMsg(w, "message", []byte("Benutzername und/oder Passwort falsch."))
        } else if u.Username == "" && redirectValue != "" {
            redirect = "/?redirect=" + redirectValue
           // setMsg(w, "message", []byte("Benutzername und/oder Passwort falsch."))
        } else {
            redirect = "/"
           // setMsg(w, "message", []byte("Benutzername und/oder Passwort falsch."))
        }
    }
    http.Redirect(w, r, redirect, 302)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    clearSession(w)
    http.Redirect(w, r, "/" , 302)
}