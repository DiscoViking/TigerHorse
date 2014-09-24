package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSON API for tigerhorse.

func StartJSON(s Storage) {
	http.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			msg := fmt.Sprintf("Cannot %v to this resource.", r.Method)
			http.Error(w, msg, 400)
			return
		}

		people, err := s.GetPeople()
		if err != nil {
			log.Print(err)
			return
		}

		data, err := json.Marshal(struct{ People []*Person }{people})
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Fprint(w, string(data))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
