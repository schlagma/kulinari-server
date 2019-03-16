package main

import (
    "html/template"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "time"
)

func NewRecipeHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    tpl := template.Must(template.ParseFiles("templates/new_recipe.html", "templates/parts.html"))
    data := RecipeData{}

    userId := 1
    data.Recipes = recipeSidebarList(userId)
    data.Categories = allCatList()
    data.URL = r.URL.Path
    data.BaseURL = "/"
    
    tpl.Execute(w, data)
}

func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
    var nId string
    var uid int

    db := dbConn()
    selDB, err := db.Query("SELECT MAX(uid) FROM recipes")
    checkErr(err)

    for selDB.Next() {
        err = selDB.Scan(&uid)
        checkErr(err)
        uid = uid + 1
        nId = strconv.Itoa(uid)
    }

    title := r.FormValue("title")
    ingredientsForValue := r.FormValue("ingredientsForValue")
    ingredientsForType := r.FormValue("ingredientsForType")
    ingredients1title := r.FormValue("ingredients1title")
    ingredients1 := r.FormValue("ingredients1")
    ingredients2title := r.FormValue("ingredients2title")
    ingredients2 := r.FormValue("ingredients2")
    ingredients3title := r.FormValue("ingredients3title")
    ingredients3 := r.FormValue("ingredients3")
    ingredients4title := r.FormValue("ingredients4title")
    ingredients4 := r.FormValue("ingredients4")
    directions := r.FormValue("directions")
    comments := r.FormValue("comments")
    interCategories := r.Form["categories"]

    var categories string
    for i := 0; i < len(interCategories); i++ {
        if i > 0 {
            categories = categories + ",[" + interCategories[i] + "]"
        } else {
            categories = "[" + interCategories[i] + "]"
        }
    }


    img := 0
    editedby := "[1]"
    editedon := time.Now().Format("2006-01-02 15:04:05")
    access := "[1]"

    insForm, err := db.Prepare("INSERT INTO recipes(uid,title,ingredients_for_value,ingredients_for_type,ingredients1,ingredients1_title,ingredients2,ingredients2_title,ingredients3,ingredients3_title,ingredients4,ingredients4_title,directions,comments,editedon,editedby,access,img,categories) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
    checkErr(err)
    _, err = insForm.Exec(nId,title,ingredientsForValue,ingredientsForType,ingredients1,ingredients1title,ingredients2,ingredients2title,ingredients3,ingredients3title,ingredients4,ingredients4title,directions,comments,editedon,editedby,access,img,categories)
    checkErr(err)
    defer insForm.Close()

    defer db.Close()

    redirect := "/recipes/view/" + nId
    http.Redirect(w, r, redirect, 302)
}