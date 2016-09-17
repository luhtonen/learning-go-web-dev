package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

const (
	Port    = ":8080"
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "cms"
	DBPass  = "cms123"
	DBDbase = "cms"
)

var database *sql.DB

type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}

func (p Page) TruncatedText() template.HTML {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 150 {
			return p.Content[:i] + ` ...`
		}
	}
	return p.Content
}

func servePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		log.Println("Couldn't get page:", pageGUID, err.Error())
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func redirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func initDB() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Could not connect to database", err.Error())
	}
	database = db
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	var Pages []Page
	pages, err := database.Query("SELECT page_title, page_content, page_date, page_guid FROM pages ORDER BY ? DESC", "page_date")
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	defer pages.Close()
	for pages.Next() {
		var thisPage Page
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

func main() {
	initDB()
	router := mux.NewRouter()
	router.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", servePage)
	router.HandleFunc("/", redirIndex)
	router.HandleFunc("/home", serveIndex)
	http.Handle("/", router)
	log.Println("Listening on port", Port)
	http.ListenAndServe(Port, nil)
}
