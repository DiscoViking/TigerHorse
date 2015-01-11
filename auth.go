package main

import (
	"encoding/base64"
	"net/http"
	"strings"
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

func validate(username, password string) bool {
	if username == "username" && password == "password" {
		return true
	}
	return false
}
