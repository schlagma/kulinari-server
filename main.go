package main

import (
    "fmt"
    "html/template"
    "net/http"
    "log"
    "os"
    "strings"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/gorilla/securecookie"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "sort"
    "encoding/json"
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
    Id       int
    Title    string
    RecImgExist bool
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
    DarkMode bool
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
    selDB2, err := db.Query("SELECT uid, title, img FROM recipes WHERE access LIKE concat('%[', ?, ']%')", userId)

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

func darkMode() bool {
    db := dbConn()
    selDB, err := db.Query("SELECT darkmode FROM settings WHERE uid=1")
    var darkMode bool
    for selDB.Next() {
        var darkmode int
        err = selDB.Scan(&darkmode)
        checkErr(err)
        darkMode = false
        if darkmode == 1 {
            darkMode = true
        }
    }
    defer db.Close()
    return darkMode
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
    r.HandleFunc("/settings/forms/darkmode", SDarkModeHandler).Methods("POST")
    r.HandleFunc("/settings/forms/colors", SColorsHandler).Methods("POST")
    r.HandleFunc("/help", HelpHandler)
    r.HandleFunc("/api/0.1/recipes/all", ApiRecListHandler)
    r.HandleFunc("/api/0.1/recipes/{id:[0-9]+}", ApiRecipeHandler)
    
    staticFileDirectory := http.Dir("./assets/")
    staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
    r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

    log.Fatal(http.ListenAndServe(":" + config.Port, r))
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
