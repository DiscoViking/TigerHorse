package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
)

func Serve(s Storage) {
	http.HandleFunc("/people/", func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error

		if r.Method != "GET" {
			msg := fmt.Sprintf("Cannot %v to this resource.", r.Method)
			http.Error(w, msg, 400)
			return
		}

		// Begin by getting all the people from the database.
		people, err := s.GetPeople()
		if err != nil {
			log.Print(err)
			return
		}

		// Check if we were asked for a specific person.
		name := path.Base(r.URL.Path)
		if name == "people" {
			// Not asked for a person.
			ServePeopleJSON(w, people)
		} else {
			// Asked for specific person.
			var p *Person

			// Find the right one.
			for _, x := range people {
				if x.Name == name {
					p = x
					break
				}
			}

			// Return their data.
			ServePersonJSON(w, p)
		}
		fmt.Fprint(w, string(data))
		return
	})

	// Handle non-API pages.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		t, err := template.ParseFiles("client/main.tmpl")
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 503)
		}
		t.Execute(w, people)
	})

	http.HandleFunc("/transaction/", func(w http.ResponseWriter, r *http.Request) {
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

		t, err := template.ParseFiles("client/transaction.tmpl")
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 503)
		}
		t.Execute(w, people)
	})

	// Serve vendor files.
	http.HandleFunc("/vendor/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			msg := fmt.Sprintf("Cannot %v to this resource.", r.Method)
			http.Error(w, msg, 400)
			return
		}

		http.ServeFile(w, r, "client/"+r.URL.Path)
	})

	// Serve data files.
	http.HandleFunc("/data/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			msg := fmt.Sprintf("Cannot %v to this resource.", r.Method)
			http.Error(w, msg, 400)
			return
		}

		http.ServeFile(w, r, "client/"+r.URL.Path)
	})

	// Serve custom css.
	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			msg := fmt.Sprintf("Cannot %v to this resource.", r.Method)
			http.Error(w, msg, 400)
			return
		}

		http.ServeFile(w, r, "client/"+r.URL.Path)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
