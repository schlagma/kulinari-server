package main

import (
    "html/template"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "strings"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "gopkg.in/russross/blackfriday.v2"
)

func ShareRecipeHandler(w http.ResponseWriter, r *http.Request) {
    tpl, _ := template.ParseFiles("templates/share_recipe.gohtml", "templates/parts.gohtml")
    
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    selDB0, err := db.Query("SELECT uid, title FROM recipes WHERE uid=?", nId)
    selDB, err := db.Query("SELECT uid, title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, img, directions, comments, ingredients_for_value, ingredients_for_type, categories, editedon, editedby FROM recipes WHERE uid=?", nId)

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
        tpl, _ = template.ParseFiles("templates/share_recipe.gohtml", "templates/parts.gohtml", "templates/notfound-recipe.gohtml")
        data.Title = "Nichts gefunden"
    }
    
    if notFound == false {
        for selDB.Next() {
            var uid, ingredients_for_value int
            var title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, directions, comments, ingredients_for_type, categories, editedon, editedby string
            var img bool
            
            err = selDB.Scan(&uid, &title, &ingredients1_title, &ingredients1, &ingredients2_title, &ingredients2, &ingredients3_title, &ingredients3, &ingredients4_title, &ingredients4, &img, &directions, &comments, &ingredients_for_value, &ingredients_for_type, &categories, &editedon, &editedby)
            
            queryValues := r.URL.Query()
            forValue := queryValues.Get("for")
            var forValInt int
            var forFactor float64
            if forValue != "" {
                
                forValInt, err = strconv.Atoi(forValue)
                forValFloat, err := strconv.ParseFloat(forValue, 64)
                if forValFloat < 1 {
                    forValFloat = 1
                    forValInt = 1
                }
                ingValFloat := float64(ingredients_for_value)
                forFactor = forValFloat / ingValFloat
                checkErr(err)
            } else {
                forValInt = ingredients_for_value
            }
            
            var interDir template.HTML
            var interDir3 string
            interDir1 := strings.Split(directions,"\n")

            for i := 0; i < len(interDir1); i++ {
                if string(interDir1[i]) != "" && string(interDir1[i]) != "\r" {
                    interDir2 := blackfriday.Run([]byte(interDir1[i]))
                    interDir3 = interDir3 + string(interDir2)
                }
            }
            tagsToRemove := []string{"div", "script", "style", "form", "datalist", "output", "input", "button", "article", "aside", "bdi", "details", "dialog", "figcaption", "figure", "footer", "header", "main", "mark", "meter", "nav", "progress", "rp", "rt", "ruby", "section", "summary", "time", "wbr"}
            for i := 0; i < len(tagsToRemove); i++ {
                interDir3 = strings.Replace(interDir3, "<" + tagsToRemove[i] + "*>", "", -1)
                interDir3 = strings.Replace(interDir3, "</" + tagsToRemove[i] + ">", "", -1)
            }
            interDir3 = strings.Replace(interDir3, "on*=\"*\"", "", -1)
            interDir3 = strings.Replace(interDir3, "style=\"*\"", "", -1)
            interDir = template.HTML(interDir3)

            var interCom template.HTML
            var interCom3 string
            interCom1 := strings.Split(comments,"\n")

            for i := 0; i < len(interCom1); i++ {
                if string(interCom1[i]) != "" && string(interCom1[i]) != "\r" {
                    interCom2 := blackfriday.Run([]byte(interCom1[i]))
                    interCom3 = interCom3 + string(interCom2)
                }
            }
            for i := 0; i < len(tagsToRemove); i++ {
                interCom3 = strings.Replace(interCom3, "<" + tagsToRemove[i] + "*>", "", -1)
                interCom3 = strings.Replace(interCom3, "</" + tagsToRemove[i] + ">", "", -1)
            }
            interCom3 = strings.Replace(interCom3, "on*=\"*\"", "", -1)
            interCom3 = strings.Replace(interCom3, "style=\"*\"", "", -1)
            interCom = template.HTML(interCom3)
            
            data.CommentsExist = !(comments=="")
            data.CategoriesExist = !(categories=="")
            data.Ing1TitleExist = !(ingredients1_title=="")
            data.Ing2TitleExist = !(ingredients2_title=="")
            data.Ing3TitleExist = !(ingredients3_title=="")
            data.Ing4TitleExist = !(ingredients4_title=="")
            
            if ingredients1 == "" {
                data.Ing1Exist = false
            } else {
                data.Ing1Exist = true
                data.Ingredients1 = ingProcess(ingredients1, forValue, forFactor)
            }
            
            if ingredients2 == "" {
                data.Ing2Exist = false
            } else {
                data.Ing2Exist = true
                data.Ingredients2 = ingProcess(ingredients2, forValue, forFactor)
            }
            
            if ingredients3 == "" {
                data.Ing3Exist = false
            } else {
                data.Ing3Exist = true
                data.Ingredients3 = ingProcess(ingredients3, forValue, forFactor)
            }
            
            if ingredients4 == "" {
                data.Ing4TitleExist = false
            } else {
                data.Ing4TitleExist = true
                data.Ingredients4 = ingProcess(ingredients4, forValue, forFactor)
            }

            interEdit, err := time.Parse("2006-01-02 15:04:05", editedon)
            checkErr(err)
            
            data.Id = uid
            data.Title = title
            data.Ingredients1Title = ingredients1_title
            data.Ingredients2Title = ingredients2_title
            data.Ingredients3Title = ingredients3_title
            data.Ingredients4Title = ingredients4_title
            data.Directions = interDir
            data.Comments = template.HTML(interCom)
            data.ImgExist = img
            data.IngredientsForValue = forValInt
            data.IngredientsForType = ingredients_for_type
            data.EditedOn = interEdit.Format("02.01.2006 15:04")
        }
    }

    data.URL = r.URL.Path
    data.BaseURL = "/"

    tpl.Execute(w, data)
    defer db.Close()
}