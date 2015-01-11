package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	authFile = ".auth"
)

var (
	username string
	password string
)

func basicAuth(pass http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if len(r.Header["Authorization"]) == 0 {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"TigerHorse\"")
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validate(pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func validate(u, p string) bool {
	if u == username && p == password {
		return true
	}
	return false
}

func loadLoginDetails() error {
	f, err := os.Open(authFile)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)

	// Read username from first line.
	if !s.Scan() {
		return s.Err()
	}

	username = s.Text()

	// Read password from second line.
	if !s.Scan() {
		return s.Err()
	}

	password = s.Text()

	log.Printf("Username: %v, Password: %v\n", username, password)

	return nil
}
