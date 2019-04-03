package main

import (
    "html/template"
    "net/http"
    "io/ioutil"
    _ "github.com/go-sql-driver/mysql"
)

func SGeneralHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    db := dbConn()
    selDB, err := db.Query("SELECT themecolor, accentcolor FROM settings WHERE uid=1")
    
    tpl := template.Must(template.ParseFiles("templates/settings_general.gohtml", "templates/parts.gohtml"))
    config := LoadConfig("config/config.json")
    data := SGeneralData{}
    data.DBhost = config.Database.Host
    data.DBname = config.Database.Name
    data.DBuser = config.Database.User
    data.Lang = config.Lang

    for selDB.Next() {
        var themecolor, accentcolor string
        err = selDB.Scan(&themecolor, &accentcolor)
        checkErr(err)
        data.ThemeColor = themecolor
        data.AccentColor = accentcolor
    }
    defer db.Close()

    data.DarkMode = darkMode()
    data.URL = r.URL.Path
    data.BaseURL = "/"
    tpl.Execute(w, data)
}

/*
func SUsersHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    userstpl := template.Must(template.ParseFiles("templates/settings_user.gohtml", "templates/parts.gohtml"))
    db := dbConn()
    selDB, err := db.Query("SELECT uid, uname, fname, lname FROM  user")
    checkErr(err)
    
    data := UserData{}
    
    var userList1 []int
    var userList2, userList3, userList4 []string
    
    for selDB2.Next() {
        var uid int
        var uname string
        var userList13 []User
        err = selDB.Scan(&uid, &uname, &fname, &lname)
        checkErr(err)
        userList1 = append(userList1, uid)
        userList2 = append(userList2, uname)
        userList2 = append(userList3, fname)
        userList2 = append(userList4, lname)
        
        for i := 1; i < len(recList1); i++ {
            userList12 := User{userList1[i], userList2[i], userList3[i], userList4[i]}
            userList13 = append(userList13, userList12)
        }
            
        data.Users = userList13
    }
    
    userstpl.Execute(w, data)
    defer db.Close()
}

func SNewUserHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    tpl := template.Must(template.ParseFiles("templates/settings_user_new.gohtml", "templates/parts.gohtml"))
    db := dbConn()
    selDB, err := db.Query("SELECT MAX(uid) FROM  user")
    checkErr(err)

    data := UserData{}

    if r.Method == "POST" {
        uname := r.FormValue("uname")
        fname := r.FormValue("fname")
        lname := r.FormValue("lname")
        password := r.FormValue("password")
        
        insForm, err := db.Prepare("INSERT INTO user(uid=?, uname=?, fname=?, lname=?, password=?) VALUES(?,?,?,?,?)")
        checkErr(err)
        insForm.Exec(uid, uname, fname, lname, password)
        http.Redirect(w, r, "/settings/user/", 302)
    }
    
    tpl.Execute(w, data)
    defer db.Close()
}
*/

func SAboutHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    tpl := template.Must(template.ParseFiles("templates/settings_about.gohtml", "templates/parts.gohtml"))
    data := RecipeData{}

    data.DarkMode = darkMode()
    data.URL = r.URL.Path
    data.BaseURL = "/"
    tpl.Execute(w, data)
}

func SAboutCopyrightHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    tpl := template.Must(template.ParseFiles("templates/settings_about_copyright.gohtml", "templates/parts.gohtml"))
    data := RecipeData{}

    data.DarkMode = darkMode()
    data.URL = r.URL.Path
    data.BaseURL = "/"
    tpl.Execute(w, data)
}

func SDarkModeHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    darkmode1 := r.FormValue("darkmode")
    darkmode := 0
    if darkmode1 != "" {
        darkmode = 1
    }

    db := dbConn()
    insForm, err := db.Prepare("UPDATE settings SET darkmode=? WHERE uid=1")
    checkErr(err)
    insForm.Exec(darkmode)
    defer insForm.Close()
    defer db.Close()

    baseURL := "/"
    redirect := baseURL + "settings/general"
    
    http.Redirect(w, r, redirect, 302)
}

func SColorsHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    themecolor := r.FormValue("themecolor")
    accentcolor := r.FormValue("accentcolor")

    db := dbConn()
    insForm, err := db.Prepare("UPDATE settings SET themecolor=?, accentcolor=? WHERE uid=1")
    checkErr(err)
    insForm.Exec(themecolor, accentcolor)
    defer insForm.Close()
    defer db.Close()

    d1 := ":root {--theme-color:" + themecolor + ";--accent-color:" + accentcolor + ";}"
    d2 := []byte(d1)
    err = ioutil.WriteFile("assets/css/colors.css", d2, 0644)
    checkErr(err)

    baseURL := "/"
    redirect := baseURL + "settings/general"
    
    http.Redirect(w, r, redirect, 302)
}

func HelpHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    abouttpl := template.Must(template.ParseFiles("templates/help.gohtml", "templates/parts.gohtml"))
    data := RecipeData{}

    data.URL = r.URL.Path
    data.BaseURL = "/"
    abouttpl.Execute(w, data)
}