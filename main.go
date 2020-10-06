package main

import (
	"net/http"
)

func main() {
	fsdb := &FileDB{"db.txt"}
	fsdb.Init()
	app := &App{fsdb}
	user := &User{}
	Creds = map[string]string{"john": "secret"}
	http.HandleFunc("/item/", user.authHandler(app.itemHandler))
	http.ListenAndServe(":8080", nil)
}
