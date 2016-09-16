package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

const (
	Port = ":8080"
)

func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	fileName := "files/" + pageID + ".html"
	http.ServeFile(w, r, fileName)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
	http.Handle("/", router)
	fmt.Println("Listening on port", Port)
	http.ListenAndServe(Port, nil)
}
