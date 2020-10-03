package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type App struct {
	db Database
}

type User struct {
	name string
	auth bool
}

func (app *App) itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// new item
	case http.MethodPost:
		defer r.Body.Close()
		item := NewItem()
		// "validation"
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// insertion
		if err := app.db.InsertItem(item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		// TODO url sanitization ?
		name := strings.Split(r.URL.Path, "/")[2]
		item, err := app.db.GetItem(name)
		// TODO inspect error
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// write json
		js, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
