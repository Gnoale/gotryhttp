package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type App struct {
	db Database
}

type User struct {
	Name     string
	Password string
}

var (
	Creds map[string]string
)

func (u *User) Validate() error {
	if pass, ok := Creds[u.Name]; ok {
		if pass == u.Password {
			return nil
		}
	}
	return errors.New("wrong username / password")
}

func (u *User) authHandler(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		u.Name, u.Password, ok = r.BasicAuth()
		err := u.Validate()
		if err != nil || !ok {
			w.Header().Set("WWW-Authenticate", "Basic realm=colombus")
			// w.WriteHeader(401)
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
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
