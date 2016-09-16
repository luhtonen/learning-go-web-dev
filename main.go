package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
	Title   string
	Content string
	Date    string
}

func servePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	thisPage := Page{}
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE id=?", pageID).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
	if err != nil {
		fmt.Println("Couldn't get page:", pageID, err.Error())
	}
	html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`
	fmt.Fprintln(w, html)
}

func initDB() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Could not connect to database", err.Error())
	}
	database = db
}

func main() {
	initDB()
	router := mux.NewRouter()
	router.HandleFunc("/pages/{id:[0-9]+}", servePage)
	http.Handle("/", router)
	fmt.Println("Listening on port", Port)
	http.ListenAndServe(Port, nil)
}
