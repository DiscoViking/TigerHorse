package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
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

func PostNewTransaction(s Storage, r *http.Request) error {
	var tx Transaction
	d := json.NewDecoder(r.Body)
	err := d.Decode(&tx)

	if err != nil {
		return errors.New("Failed to parse JSON body.")
	}

	// Stamp with current time.
	tx.Time = time.Now()

	err = s.AddTransaction(&tx)

	if err != nil {
		return errors.New("Error adding transaction to database.")
	}

	return nil
}
