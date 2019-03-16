package main

import (
    "net/http"
	_ "github.com/go-sql-driver/mysql"
    "strings"
    "strconv"
	"sort"
    "encoding/json"
    "time"
    "github.com/gorilla/mux"
)

func ApiRecListHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    db := dbConn()
    data := RecList{}
    selDB, err := db.Query("SELECT uid, title FROM recipes")
    checkErr(err)

    var recList1 []int
    var recList2 []string
    var recList13 []Recipe
    
    for selDB.Next() {
        var uid int
        var title string
        err = selDB.Scan(&uid, &title)
        checkErr(err)
        recList1 = append(recList1, uid)
        recList2 = append(recList2, title)
        
        for i := 0; i < len(recList1); i++ {
            recList12 := Recipe{recList1[i], recList2[i]}
            recList13 = append(recList13, recList12)
        }
        
        recList13 = rmDuplRec(recList13)
        
        sort.Sort(RecipesByTitle(recList13))
            
        data.Recipes = recList13
    }

    m, _ := json.Marshal(data)
    w.Write(m)
}

func ApiRecipeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    selDB, err := db.Query("SELECT uid, title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, img, directions, comments, ingredients_for_value, ingredients_for_type, categories, editedon FROM recipes WHERE uid=?", nId)
    selDB2, err := db.Query("SELECT uid, title FROM categories")
    
    checkErr(err)
    
    data := ApiRecipeData {}
    
    var catList3 []string
    
    for selDB.Next() {
        var uid, ingredients_for_value int
        var title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, directions, comments, ingredients_for_type, categories, editedon string
        var img bool
        
        err = selDB.Scan(&uid, &title, &ingredients1_title, &ingredients1, &ingredients2_title, &ingredients2, &ingredients3_title, &ingredients3, &ingredients4_title, &ingredients4, &img, &directions, &comments, &ingredients_for_value, &ingredients_for_type, &categories, &editedon)
        
        queryValues := r.URL.Query()
        forValue := queryValues.Get("for")
        var forValInt int
        var forFactor float64
        if forValue != "" {
            forValInt, err = strconv.Atoi(forValue)
            forValFloat, err := strconv.ParseFloat(forValue, 64)
            ingValFloat := float64(ingredients_for_value)
            forFactor = forValFloat / ingValFloat
            checkErr(err)
        } else {
            forValInt = ingredients_for_value
        }
        
        directions = strings.Replace(directions, "\r", "", -1)
        interDir := strings.Split(directions,"\n")
        comments = strings.Replace(comments, "\r", "", -1)
        interCom := strings.Split(comments,"\n")
        
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
        
        catList31 := strings.Replace(categories, "[", "", -1)
        catList31 = strings.Replace(catList31, "]", "", -1)
        catList3 = strings.Split(catList31,",")
        
        interEdit, err := time.Parse("2006-01-02 15:04:05", editedon)
        checkErr(err)
        
        data.Id = uid
        data.Title = title
        data.Ingredients1Title = ingredients1_title
        data.Ingredients2Title = ingredients2_title
        data.Ingredients3Title = ingredients3_title
        data.Ingredients4Title = ingredients4_title
        data.Directions = interDir
        data.Comments = interCom
        data.ImgExist = img
        data.IngredientsForValue = forValInt
        data.IngredientsForType = ingredients_for_type
        data.EditedOn = interEdit.Format("2006-01-02 15:04:05")
    }
    
    var catList1 []int
    var catList2 []string
    var catList13, catList14 []Category
    
    if data.CategoriesExist == true {
        for selDB2.Next() {
            var uid int
            var title string
            err = selDB2.Scan(&uid, &title)
            checkErr(err)
            catList1 = append(catList1, uid)
            catList2 = append(catList2, title)
            
            for i := 0; i < len(catList1); i++ {
                catList12 := Category{catList1[i], catList2[i]}
                catList13 = append(catList13, catList12)
            }
        }
        
        for i := 0; i < len(catList3); i++ {
            for j := 0; j < len(catList13); j++ {
                k, err := strconv.Atoi(catList3[i])
                checkErr(err)
                if k == catList13[j].Id {
                    catList14 = append(catList14, catList13[j])
                }
            }
        }
        
        catList14 = rmDuplCat(catList14)
        
        sort.Sort(CategoriesByTitle(catList14))
            
        data.Categories = catList14
    }

    m, _ := json.Marshal(data)
    w.Write(m)
}
