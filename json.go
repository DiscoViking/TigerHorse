package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSON API for tigerhorse.

func ServePeopleJSON(w http.ResponseWriter, people []*Person) {
	data, err := json.Marshal(struct{ People []*Person }{people})
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 503)
		return
	}
	fmt.Fprint(w, string(data))
	return
}

func ServePersonJSON(w http.ResponseWriter, p *Person, txs []*Transaction) {
	// TODO: Return all transaction data.
	data, err := json.Marshal(struct {
		Person       *Person
		Transactions []*Transaction
	}{p, txs})
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 503)
		return
	}
	fmt.Fprint(w, string(data))
	return
}
