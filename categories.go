package main

import (
    "html/template"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "strings"
    "sort"
    "github.com/gorilla/mux"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    tpl := template.Must(template.ParseFiles("templates/categories.gohtml", "templates/parts.gohtml"))
    data := RecipeData {}

    userId := 1
    data.Recipes = recipeSidebarList(userId)
    data.Categories = allCatList()
    data.Title = "Kategorien"
    data.URL = r.URL.Path
    data.BaseURL = "/"
    
    tpl.Execute(w, data)
}

func CatContentHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    cattpl := template.Must(template.ParseFiles("templates/categories_content.gohtml", "templates/parts.gohtml"))
    
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    selDB2, err := db.Query("SELECT uid, title, img, categories FROM recipes")
    selDB3, err := db.Query("SELECT title FROM categories WHERE uid=?", nId)
    
    checkErr(err)
        
    data := RecipeData {}
    
    var recList1 []int
    var recList2, recListCat1 []string
    var recList13, recListCat3 []Recipe
    var recList3 []bool
    
    for selDB2.Next() {
        var uid int
        var title, categories string
        var img bool
        err = selDB2.Scan(&uid, &title, &img, &categories)
        checkErr(err)
        recList1 = append(recList1, uid)
        recList2 = append(recList2, title)
        recList3 = append(recList3, img)
        
        for i := 1; i < len(recList1); i++ {
            recList12 := Recipe{recList1[i], recList2[i], recList3[i]}
            recList13 = append(recList13, recList12)
        }
    
        recList13 = rmDuplRec(recList13)
        
        sort.Sort(RecipesByTitle(recList13))
            
        data.Recipes = recList13
        
        
        recListCat1 = append(recListCat1, categories)
        
        for i := 0; i < len(recListCat1); i++ {
            if strings.Contains(recListCat1[i], "[" + nId + "]") {
                recListCat2 := Recipe{recList1[i], recList2[i], recList3[i]}
                recListCat3 = append(recListCat3, recListCat2)
            }
        }
        
        recListCat3 = rmDuplRec(recListCat3)
        sort.Sort(RecipesByTitle(recListCat3))
            
        data.CatContent = recListCat3
    }
    
    for selDB3.Next() {
        var title string
        err = selDB3.Scan(&title)
        checkErr(err)
        data.Title = title
    }

    data.URL = r.URL.Path
    data.BaseURL = "/"
    
    cattpl.Execute(w, data)
    defer db.Close()
}