package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"time"
)

const (
	Port = ":8080"
)

func serveDynamic(w http.ResponseWriter, r *http.Request) {
	response := "The time is now " + time.Now().String()
	fmt.Fprintln(w, response)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static.html")
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/test", TestHandler)
	router.HandleFunc("/", serveDynamic)
	router.HandleFunc("/static", serveStatic)
	http.Handle("/", router)
	fmt.Println("Listening on port", Port)
	http.ListenAndServe(Port, nil)
}
