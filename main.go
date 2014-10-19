package main

import "log"

func main() {
	DSN := "tigerhorse:whoopee@tcp(localhost:3306)/tigerhorse?parseTime=true&loc=GMT"

	log.Print("Establishing connection to database...")
	s, err := New(DSN)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting TigerHorse Server...")
	Serve(s)
	return
}
