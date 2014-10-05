package main

import "log"

func main() {
	DSN := "tigerhorse@tcp(localhost:3306)/tigerhorsetest?parseTime=true&loc=GMT"
	s, err := New(DSN)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting TigerHorse Server...")
	Serve(s)
	return
}
