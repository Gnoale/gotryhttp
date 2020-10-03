package main

import (
	"net/http"
)

func main() {
	fsdb := &FileDB{"db.txt"}
	fsdb.Init()
	app := &App{fsdb}
	http.HandleFunc("/item/", app.itemHandler)
	http.ListenAndServe(":8080", nil)
}
