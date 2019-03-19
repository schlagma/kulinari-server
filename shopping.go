package main

import (
    "html/template"
    "net/http"
    "sort"
    _ "github.com/go-sql-driver/mysql"
)

func ShoppingHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    abouttpl := template.Must(template.ParseFiles("templates/shopping.gohtml", "templates/parts.gohtml"))
    
    db := dbConn()
    selDB2, err := db.Query("SELECT uid, title, img FROM  recipes")

    checkErr(err)
    
    data := RecipeData{}
    
    var recList1 []int
    var recList2 []string
    var recList3 []bool
    var recList13 []Recipe
    
    for selDB2.Next() {
        var uid int
        var title string
        var img bool
        err = selDB2.Scan(&uid, &title, &img)
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
    }
    defer db.Close()

    data.URL = r.URL.Path
    
    abouttpl.Execute(w, data)
}