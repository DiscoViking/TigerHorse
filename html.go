package main

import (
	"html/template"
	"net/http"
)

// Functions for serving HTML content.

func ServePeopleHtml(w http.ResponseWriter, people []*Person) {
}

func ServePersonHtml(w http.ResponseWriter, person *Person, txs []*Transaction) {
	funcs := template.FuncMap{
		"pounds": penniesToPounds,
	}

	t := template.Must(template.New("person.tmpl").Funcs(funcs).ParseFiles("client/person.tmpl"))

	t.Execute(w, struct {
		Person       *Person
		Transactions []*Transaction
	}{person, txs})
}
