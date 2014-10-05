package main

import (
	"log"
	"net/http"
	"text/template"
)

// Functions for serving HTML content.

func ServePeopleHtml(w http.ResponseWriter, people []*Person) {
}

func ServePersonHtml(w http.ResponseWriter, person *Person, txs []*Transaction) {
	t, err := template.ParseFiles("client/person.tmpl")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 503)
	}

	t.Execute(w, struct {
		Person       *Person
		Transactions []*Transaction
	}{person, txs})
}
