package main

import (
    "html/template"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "strings"
    "strconv"
    "sort"
    "time"
)

func EditRecipeHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    tpl, _ := template.ParseFiles("templates/edit_recipe.gohtml", "templates/parts.gohtml")
    
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    selDB0, err := db.Query("SELECT uid, title FROM recipes WHERE uid=?", nId)
    selDB, err := db.Query("SELECT uid, title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, img, directions, comments, ingredients_for_value, ingredients_for_type, categories FROM recipes WHERE uid=?", nId)
    selDB3, err := db.Query("SELECT * FROM categories")
    selDB4, err := db.Query("SELECT uid, categories FROM recipes WHERE uid=?", nId)
    checkErr(err)
    
    data := RecipeData {}

    notFound := true
    
    for selDB0.Next() {
        var uid int
        var title string
        err = selDB.Scan(&uid, &title)

        if title == "" {
            notFound = false
        }
    }

    if notFound == true {
        tpl, _ = template.ParseFiles("templates/recipe.gohtml", "templates/parts.gohtml", "templates/notfound-recipe.gohtml")
        data.Title = "Nichts gefunden"
    }

    if notFound == false {
        for selDB.Next() {
            var uid, ingredients_for_value int
            var title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, directions, comments, ingredients_for_type, categories string
            var img bool
            
            err = selDB.Scan(&uid, &title, &ingredients1_title, &ingredients1, &ingredients2_title, &ingredients2, &ingredients3_title, &ingredients3, &ingredients4_title, &ingredients4, &img, &directions, &comments, &ingredients_for_value, &ingredients_for_type, &categories)
            
            checkErr(err)
            
            data.Id = uid
            data.Title = title
            data.Ingredients1Title = ingredients1_title
            data.Ingredients1Raw = ingredients1
            data.Ingredients2Title = ingredients2_title
            data.Ingredients2Raw = ingredients2
            data.Ingredients3Title = ingredients3_title
            data.Ingredients3Raw = ingredients3
            data.Ingredients4Title = ingredients4_title
            data.Ingredients4Raw = ingredients4
            data.DirectionsRaw = directions
            data.CommentsRaw = comments
            data.ImgExist = img
            data.IngredientsForValue = ingredients_for_value
            data.IngredientsForType = ingredients_for_type
        }

        var catList1 []int
        var catList2 []string
        var catList3 string
        var catList13 []Category
        var catList4 []CategorySelected
        
        for selDB3.Next() {
            var uid int
            var title string
            err = selDB3.Scan(&uid, &title)
            checkErr(err)
            catList1 = append(catList1, uid)
            catList2 = append(catList2, title)
            for i := 0; i < len(catList1); i++ {
                catList12 := Category{catList1[i], catList2[i]}
                catList13 = append(catList13, catList12)
            }
        }
        catList13 = rmDuplCat(catList13)

        for selDB4.Next() {
            var categories string
            err = selDB.Scan(&categories)
            catList3 = categories
        }

        for i := 0; i < len(catList1); i++ {
            if strings.Contains(catList3, "[" + strconv.Itoa(catList13[i].Id) + "]") {
                catList14 := CategorySelected{catList13[i].Id, catList13[i].Title, true}
                catList4 = append(catList4, catList14)
            } else {
                catList14 := CategorySelected{catList13[i].Id, catList13[i].Title, false}
                catList4 = append(catList4, catList14)
            }
        }
        sort.Sort(CategoriesSelectedByTitle(catList4))

        data.CategoriesSelected = catList4
    }
    userId := 1
    data.Recipes = recipeSidebarList(userId)
    
    if r.Method == "POST" {
        uid, err := strconv.Atoi(nId)
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
        editedon := time.Now().Format("2006-01-02 15:04:05")
        interCategories := r.Form["categories"]
        var categories string
        for i := 0; i < len(interCategories); i++ {
            if i > 0 {
                categories = categories + ",[" + interCategories[i] + "]"
            } else {
                categories = "[" + interCategories[i] + "]"
            }
        }

        insForm, err := db.Prepare("UPDATE recipes SET title=?, ingredients_for_value=?, ingredients_for_type=?, ingredients1_title=?, ingredients1=?, ingredients2_title=?, ingredients2=?, ingredients3_title=?, ingredients3=?, ingredients4_title=?, ingredients4=?, directions=?, comments=?, editedon=?, categories=? WHERE uid=?")
        checkErr(err)
        insForm.Exec(title, ingredientsForValue, ingredientsForType, ingredients1title, ingredients1, ingredients2title, ingredients2, ingredients3title, ingredients3, ingredients4title, ingredients4, directions, comments, editedon, categories, uid)
        http.Redirect(w, r, "/recipes/view/" + nId, 302)
    }
    defer db.Close()

    data.DarkMode = darkMode()
    data.URL = r.URL.Path
    data.BaseURL = "/"
    
    tpl.Execute(w, data)
}