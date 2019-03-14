package main

import (
    "fmt"
    "html/template"
    "net/http"
    "log"
    "os"
    "io/ioutil"
    "strings"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/gorilla/securecookie"
    "database/sql"
     _ "github.com/go-sql-driver/mysql"
    "sort"
    "time"
    "encoding/json"
    "gopkg.in/russross/blackfriday.v2"
    /*"golang.org/x/crypto/bcrypt"*/
)

type RecipeData struct {
    Id                      int
    Title                   string
    Ingredients1Title       string
    Ingredients1            [][]string
    Ingredients1Raw         string
    Ingredients2Title       string
    Ingredients2            [][]string
    Ingredients2Raw         string
    Ingredients3Title       string
    Ingredients3            [][]string
    Ingredients3Raw         string
    Ingredients4Title       string
    Ingredients4            [][]string
    Ingredients4Raw         string
    Directions              template.HTML
    DirectionsRaw           string
    Comments                template.HTML
    CommentsRaw             string
    Categories              []Category
    CategoriesSelected      []CategorySelected
    IngredientsForValue     int
    IngredientsForType      string
    EditedOn                string
    EditedBy                string
    CommentsExist           bool
    CategoriesExist         bool
    Ing1TitleExist          bool
    Ing2TitleExist          bool
    Ing3TitleExist          bool
    Ing4TitleExist          bool
    Ing1Exist               bool
    Ing2Exist               bool
    Ing3Exist               bool
    Ing4Exist               bool
    ImgExist                bool
    Recipes                 []Recipe
    CatContent              []Recipe
    URL                     string
    BaseURL                 string
    DarkMode                bool
}

type ApiRecipeData struct {
    Id                      int
    Title                   string
    Ingredients1Title       string
    Ingredients1            [][]string
    Ingredients2Title       string
    Ingredients2            [][]string
    Ingredients3Title       string
    Ingredients3            [][]string
    Ingredients4Title       string
    Ingredients4            [][]string
    Directions              []string
    Comments                []string
    Categories              []Category
    IngredientsForValue     int
    IngredientsForType      string
    EditedOn                string
    EditedBy                string
    CommentsExist           bool
    CategoriesExist         bool
    Ing1TitleExist          bool
    Ing2TitleExist          bool
    Ing3TitleExist          bool
    Ing4TitleExist          bool
    Ing1Exist               bool
    Ing2Exist               bool
    Ing3Exist               bool
    Ing4Exist               bool
    ImgExist                bool
}

type RecList struct {
    Recipes []Recipe
}

type Recipe struct {
    Id      int
    Title   string
}

type Category struct {
    Id      int
    Title   string
}

type CategorySelected struct {
    Id          int
    Title       string
    IsSelected  bool
}

type Ingredient struct {
    Value   int
    Unit    string
    Food    string
}

type RecipesByTitle []Recipe
type CategoriesByTitle []Category
type CategoriesSelectedByTitle []CategorySelected

type LoginData struct {
    User     string
    Redirect string
    BaseURL  string
}

type Config struct {
    Database struct {
        Host        string  `json:"host"`
        Name        string  `json:"name"`
        User        string  `json:"user"`
        Password    string  `json:"password"`
    } `json:"database"`
    Port    string  `json:"port"`
    Lang    string  `json:"lang"`
    URL     string
}

type SGeneralData struct {
    DBhost      string
    DBname      string
    DBuser      string
    Lang        string
    URL         string
    BaseURL     string
    ThemeColor  string
    AccentColor string
    DarkMode    bool
}

func dbConn() (db *sql.DB) {
    config := LoadConfig("config/config.json")
    dbDriver := "mysql"
    dbUser := config.Database.User
    dbPass := config.Database.Password
    dbName := config.Database.Name
    var dbHost string
    if config.Database.Host == "localhost" {
        dbHost = ""
    } else {
        dbHost = "tcp("+config.Database.Host+")"
    }
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbHost+"/"+dbName)
    checkErr(err)
    return db
}

func (a RecipesByTitle) Len() int           { return len(a) }
func (a RecipesByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }
func (a RecipesByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a CategoriesByTitle) Len() int           { return len(a) }
func (a CategoriesByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }
func (a CategoriesByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a CategoriesSelectedByTitle) Len() int           { return len(a) }
func (a CategoriesSelectedByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }
func (a CategoriesSelectedByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func rmDuplRec(elements []Recipe) []Recipe {
    encountered := map[Recipe]bool{}
    result := []Recipe{}

    for v := range elements {
        if encountered[elements[v]] == true {
        } else {
            encountered[elements[v]] = true
            result = append(result, elements[v])
        }
    }
    return result
}

func rmDuplCat(elements []Category) []Category {
    encountered := map[Category]bool{}
    result := []Category{}

    for v := range elements {
        if encountered[elements[v]] == true {
        } else {
            encountered[elements[v]] = true
            result = append(result, elements[v])
        }
    }
    return result
}

func rmDuplCatSelected(elements []CategorySelected) []CategorySelected {
    encountered := map[CategorySelected]bool{}
    result := []CategorySelected{}

    for v := range elements {
        if encountered[elements[v]] == true {
        } else {
            encountered[elements[v]] = true
            result = append(result, elements[v])
        }
    }
    return result
}

func ingProcess(procIng string, forValue string, forFactor float64) (interIng3 [][]string) {
    var interIng2 []string
    procIng = strings.Replace(procIng, "\r", "", -1)
    interIng1 := strings.Split(procIng,"\n")
    for i := 0; i < len(interIng1); i++ {
        interIng2 = strings.Split(interIng1[i],"§§")
        interIng3 = append(interIng3, interIng2)
    }
    if forValue != "" {
        for i := 0; i < len(interIng3); i++ {
            interIngReplace := strings.Replace(interIng3[i][0], ",", ".", -1)
            if _, err := strconv.ParseFloat(interIngReplace, 64); err == nil {
                interIngFloat, err := strconv.ParseFloat(interIngReplace, 64)
                interIngCalc := interIngFloat * forFactor
                var interIngString string
                if interIngCalc == float64(int(interIngCalc)) {
                    interIngString = strconv.Itoa(int(interIngCalc))
                } else {
                    interIngString = strings.Replace(strings.TrimRight(fmt.Sprintf("%.2f", interIngCalc), "0"), ".", ",", -1)
                }
                interIng3[i][0] = interIngString
                checkErr(err)
            }
        }
    }
    return interIng3
}

func LoadConfig(file string) Config {
    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) {
    username := getUserName(r)
    if username == "" {
        http.Redirect(w, r, "/?redirect=" + r.URL.Path , 302)
    }
}

var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))

/*
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
*/

func sliceContainsInt(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func recipeSidebarList(userId int) []Recipe {
    db := dbConn()
    selDB2, err := db.Query("SELECT uid, title FROM recipes WHERE access LIKE concat('%[', ?, ']%')", userId)

    var recList1 []int
    var recList2 []string
    var recList13 []Recipe
    
    for selDB2.Next() {
        var uid int
        var title string
        err = selDB2.Scan(&uid, &title)
        checkErr(err)
        recList1 = append(recList1, uid)
        recList2 = append(recList2, title)
        for i := 1; i < len(recList1); i++ {
            recList12 := Recipe{recList1[i], recList2[i]}
            recList13 = append(recList13, recList12)
        }
        recList13 = rmDuplRec(recList13)
        sort.Sort(RecipesByTitle(recList13)) 
    }
    defer db.Close()
    return recList13
}

func allCatList() []Category {
    db := dbConn()
    selDB3, err := db.Query("SELECT * FROM categories")
    checkErr(err)

    var catList1 []int
    var catList2 []string
    var catList13 []Category
    
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
        catList13 = rmDuplCat(catList13)
        sort.Sort(CategoriesByTitle(catList13))
    }
    defer db.Close()
    return catList13
}

func main() {
    config := LoadConfig("config/config.json")

    r := mux.NewRouter()
    r.HandleFunc("/", StartHandler)
    r.HandleFunc("/login", LoginHandler).Methods("POST")
    r.HandleFunc("/logout", LogoutHandler).Methods("POST")
    r.HandleFunc("/recipes/new", NewRecipeHandler)
    r.HandleFunc("/recipes/create", CreateRecipeHandler).Methods("POST")
    r.HandleFunc("/recipes/view/{id:[0-9]+}", RecipeHandler)
    r.HandleFunc("/recipes/edit/{id:[0-9]+}", EditRecipeHandler)
    r.HandleFunc("/recipes/delete/{id:[0-9]+}", DeleteRecipeHandler)
    r.HandleFunc("/recipes/cat/{id:[0-9]+}", CatContentHandler)
    //r.HandleFunc("/recipes/export/{id:[0-9]+}", ExportRecipeHandler)
    r.HandleFunc("/recipes/s/{id:[0-9]+}", ShareRecipeHandler)
    r.HandleFunc("/recipes", CategoriesHandler)
    r.HandleFunc("/shopping", ShoppingHandler)
   // r.HandleFunc("/settings/users/new", SNewUserHandler)
    r.HandleFunc("/settings/general", SGeneralHandler)
   // r.HandleFunc("/settings/users", SUsersHandler)
    r.HandleFunc("/settings/about", SAboutHandler)
    r.HandleFunc("/settings/about/copyright", SAboutCopyrightHandler)
    r.HandleFunc("/settings/forms/colors", SColorsHandler).Methods("POST")
    r.HandleFunc("/help", HelpHandler)
    r.HandleFunc("/api/0.1/recipes/all", ApiRecListHandler)
    r.HandleFunc("/api/0.1/recipes/{id:[0-9]+}", ApiRecipeHandler)
    
    staticFileDirectory := http.Dir("./assets/")
    staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
    r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

    log.Fatal(http.ListenAndServe(":" + config.Port, r))
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
    tpl, _ := template.ParseFiles("templates/login.html", "templates/parts.html")
    data := LoginData{}

    queryValues := r.URL.Query()
    userValue := queryValues.Get("user")
    redirectValue := queryValues.Get("redirect")
    data.User = userValue
    data.Redirect = redirectValue

    username := getUserName(r)
    if username != "" && redirectValue == "" {
        http.Redirect(w, r, "/recipes" , 302)
    } else if username != "" && redirectValue != "" {
        http.Redirect(w, r, redirectValue, 302)
    }
    
    data.BaseURL = "/"

    tpl.Execute(w, data)
}

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    tpl, _ := template.ParseFiles("templates/recipe.html", "templates/parts.html")
    
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    userId := 1
    selDB0, err := db.Query("SELECT uid, title FROM recipes WHERE uid=?", nId)
    selDB, err := db.Query("SELECT uid, title, ingredients1_title, ingredients1, ingredients2_title, ingredients2, ingredients3_title, ingredients3, ingredients4_title, ingredients4, img, directions, comments, ingredients_for_value, ingredients_for_type, categories, editedon, editedby FROM recipes WHERE uid=?", nId)
    selDB3, err := db.Query("SELECT uid, title FROM categories")

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
        tpl, _ = template.ParseFiles("templates/recipe.html", "templates/parts.html", "templates/notfound-recipe.html")
        data.Title = "Nichts gefunden"
    }

    var catList3 []string
    var editedById int
    
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
            
            catList31 := strings.Replace(categories, "[", "", -1)
            catList31 = strings.Replace(catList31, "]", "", -1)
            catList3 = strings.Split(catList31,",")

            editedById1 := strings.Replace(editedby, "[", "", -1)
            editedById1 = strings.Replace(editedById1, "]", "", -1)
            editedById, err = strconv.Atoi(editedById1)
            checkErr(err)

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

        editedById = 1

        /*selDB4, err := db.Query("SELECT fname, lname FROM user WHERE uid=?", editedById)
        checkErr(err)

        for selDB4.Next() {
            var fname, lname string
            err = selDB.Scan(&fname, &lname)
            checkErr(err)
            editedByName := fname + " " + lname
            data.EditedBy = editedByName
        }*/

        data.EditedBy = strconv.Itoa(editedById)
    }

    data.Recipes = recipeSidebarList(userId)

    var catList1 []int
    var catList2 []string
    var catList13, catList14 []Category
    
    if notFound == false {
        if data.CategoriesExist == true {
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
    }

    data.URL = r.URL.Path
    data.BaseURL = "/"

    tpl.Execute(w, data)
    defer db.Close()
}

func EditRecipeHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    tpl, _ := template.ParseFiles("templates/edit_recipe.html", "templates/parts.html")
    
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
        tpl, _ = template.ParseFiles("templates/recipe.html", "templates/parts.html", "templates/notfound-recipe.html")
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

    data.URL = r.URL.Path
    data.BaseURL = "/"
    
    tpl.Execute(w, data)
}

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

func ShareRecipeHandler(w http.ResponseWriter, r *http.Request) {
    tpl, _ := template.ParseFiles("templates/share_recipe.html", "templates/parts.html")
    
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
        tpl, _ = template.ParseFiles("templates/share_recipe.html", "templates/parts.html", "templates/notfound-recipe.html")
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

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    tpl := template.Must(template.ParseFiles("templates/categories.html", "templates/parts.html"))
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
    
    cattpl := template.Must(template.ParseFiles("templates/categories_content.html", "templates/parts.html"))
    
    db := dbConn()
    params := mux.Vars(r)
    nId := params["id"]
    selDB2, err := db.Query("SELECT uid, title, categories FROM recipes")
    selDB3, err := db.Query("SELECT title FROM categories WHERE uid=?", nId)
    
    checkErr(err)
        
    data := RecipeData {}
    
    var recList1 []int
    var recList2, recListCat1 []string
    var recList13, recListCat3 []Recipe
    
    for selDB2.Next() {
        var uid int
        var title, categories string
        err = selDB2.Scan(&uid, &title, &categories)
        checkErr(err)
        recList1 = append(recList1, uid)
        recList2 = append(recList2, title)
        
        for i := 1; i < len(recList1); i++ {
            recList12 := Recipe{recList1[i], recList2[i]}
            recList13 = append(recList13, recList12)
        }
    
        recList13 = rmDuplRec(recList13)
        
        sort.Sort(RecipesByTitle(recList13))
            
        data.Recipes = recList13
        
        
        recListCat1 = append(recListCat1, categories)
        
        for i := 0; i < len(recListCat1); i++ {
            if strings.Contains(recListCat1[i], "[" + nId + "]") {
                recListCat2 := Recipe{recList1[i], recList2[i]}
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

func ShoppingHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    abouttpl := template.Must(template.ParseFiles("templates/shopping.html", "templates/parts.html"))
    
    db := dbConn()
    selDB2, err := db.Query("SELECT uid, title FROM  recipes")

    checkErr(err)
    
    data := RecipeData{}
    
    var recList1 []int
    var recList2 []string
    var recList13 []Recipe
    
    for selDB2.Next() {
        var uid int
        var title string
        err = selDB2.Scan(&uid, &title)
        checkErr(err)
        recList1 = append(recList1, uid)
        recList2 = append(recList2, title)
        
        for i := 1; i < len(recList1); i++ {
            recList12 := Recipe{recList1[i], recList2[i]}
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

func SGeneralHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)

    db := dbConn()
    selDB, err := db.Query("SELECT themecolor, accentcolor, darkmode FROM settings WHERE uid=1")
    
    tpl := template.Must(template.ParseFiles("templates/settings_general.html", "templates/parts.html"))
    config := LoadConfig("config/config.json")
    data := SGeneralData{}
    data.DBhost = config.Database.Host
    data.DBname = config.Database.Name
    data.DBuser = config.Database.User
    data.Lang = config.Lang

    for selDB.Next() {
        var darkmode int
        var themecolor, accentcolor string
        err = selDB.Scan(&themecolor, &accentcolor, &darkmode)
        checkErr(err)
        data.ThemeColor = themecolor
        data.AccentColor = accentcolor
        data.DarkMode = false
        if darkmode == 1 {
            data.DarkMode = true
        }
    }
    defer db.Close()

    data.URL = r.URL.Path
    data.BaseURL = "/"
    tpl.Execute(w, data)
}

/*
func SUsersHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    userstpl := template.Must(template.ParseFiles("templates/settings_user.html", "templates/parts.html"))
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
    
    tpl := template.Must(template.ParseFiles("templates/settings_user_new.html", "templates/parts.html"))
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
    
    tpl := template.Must(template.ParseFiles("templates/settings_about.html", "templates/parts.html"))
    data := RecipeData{}

    data.URL = r.URL.Path
    data.BaseURL = "/"
    tpl.Execute(w, data)
}

func SAboutCopyrightHandler(w http.ResponseWriter, r *http.Request) {
    isLoggedIn(w,r)
    
    tpl := template.Must(template.ParseFiles("templates/settings_about_copyright.html", "templates/parts.html"))
    data := RecipeData{}

    data.URL = r.URL.Path
    data.BaseURL = "/"
    tpl.Execute(w, data)
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
    
    abouttpl := template.Must(template.ParseFiles("templates/help.html", "templates/parts.html"))
    data := RecipeData{}

    data.URL = r.URL.Path
    data.BaseURL = "/"
    abouttpl.Execute(w, data)
}

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



func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
