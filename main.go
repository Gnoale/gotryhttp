package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fsdb := &FileDB{"db.txt"}
	fsdb.Init()
	app := &App{fsdb}
	//user := &User{}
	//Creds = map[string]string{"john": "secret"}
	r := mux.NewRouter()
	r.HandleFunc("/item", app.newItem).Methods("POST")
	r.HandleFunc("/item/{name}", app.getItem).Methods("GET")

	http.ListenAndServe(":8080", r)
}
