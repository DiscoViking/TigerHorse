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

func ServePersonJSON(w http.ResponseWriter, p *Person) {
	// TODO: Return all transaction data.
	data, err := json.Marshal(struct{ Person *Person }{p})
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 503)
		return
	}
	fmt.Fprint(w, string(data))
	return
}
