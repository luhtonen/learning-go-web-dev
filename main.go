package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	if len(p.Content) > 150 {
		return p.Content[:150] + ` ...`
	}
	return p.Content
}

type JSONResponse struct {
	Fields map[string]string
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

func apiPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	var thisPage Page
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	_, err = json.Marshal(thisPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, thisPage)
}

func apiCommentPost(w http.ResponseWriter, r *http.Request) {
	var commentAdded bool
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")

	res, err := database.Exec("INSERT INTO comments SET comment_name=?, comment_email=?, comment_text=?", name, email, comments)

	if err != nil {
		log.Println(err.Error())
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		commentAdded = false
	} else {
		commentAdded = true
	}
	commentAddedBool := strconv.FormatBool(commentAdded)
	var resp JSONResponse
	resp.Fields["id"] = string(id)
	resp.Fields["added"] =  commentAddedBool
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonResp)
}

func main() {
	initDB()
	routes := mux.NewRouter()
	routes.HandleFunc("/api/pages", apiPage).Methods("GET")
	routes.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", apiPage).Methods("GET")
	routes.HandleFunc("/api/comments", apiCommentPost).Methods("POST")
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", servePage)
	routes.HandleFunc("/", redirIndex)
	routes.HandleFunc("/home", serveIndex)
	http.Handle("/", routes)
	certificates, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Fatalln("Cannot load certificates:", err.Error())
	}
	tlsConfig := tls.Config{Certificates: []tls.Certificate{certificates}}
	log.Println("TLS Config", tlsConfig)
	log.Println("Listening on port", Port)
	//_, err = tls.Listen("tcp", Port, &tlsConfig)
	err = http.ListenAndServeTLS(Port, "server.pem", "server.key", nil)
	if err != nil {
		log.Fatalln("failed to listen:", err.Error())
	}
}
