package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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

			if p == nil {
				// Didn't find anyone.
				return
			}

			// Return their data.
			//ServePersonJSON(w, p)
			txs, err := p.Transactions(s)
			if err != nil {
				log.Print(err)
				return
			}

			// Hack with the transactions so the ones we didn't buy show the right value.
			// TODO: Do this properly.
			for _, tx := range txs {
				if tx.Buyer != p.Id {
					tx.Value = -(tx.Value / (int64(len(tx.Involved)) + tx.Guests))
				}
			}

			// Switch on accepted content types.
			accepts := r.Header["Accept"]

			// Html if unspecified.
			var contentType string
			if len(accepts) == 0 {
				contentType = "text/html"
			} else {
				contentTypes := strings.Split(accepts[0], ",")
				contentType = contentTypes[0]
			}

			log.Printf("Serving content-type: %s", contentType)

			switch contentType {
			case "text/html":
				ServePersonHtml(w, p, txs)
			case "text/json":
				ServePersonJSON(w, p, txs)
			}
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

		funcs := template.FuncMap{
			"pounds": penniesToPounds,
		}

		t := template.Must(template.New("main.tmpl").Funcs(funcs).ParseFiles("client/main.tmpl"))

		t.Execute(w, people)
	})

	http.HandleFunc("/transaction/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
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

		case "POST":
			fmt.Print("Received transaction POST request")
			err := PostNewTransaction(s, r)
			if err != nil {
				fmt.Print("Error adding transaction: ", err.Error())
				http.Error(w, err.Error(), 503)
			}

		default:
			msg := fmt.Sprintf("Cannot %v to this resource.", r.Method)
			http.Error(w, msg, 400)
		}

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
