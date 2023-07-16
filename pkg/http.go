package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpInit() {
	fmt.Println("Loading HTTP...")
	http.HandleFunc("/", handleIndex)

	// HTTP should pin us.
	log.Fatal(http.ListenAndServe("127.0.0.1:8087", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	jsonString := "{test: 'test'}"
	fprintf, err := fmt.Fprintf(w, "%s", string(jsonString))
	fmt.Println(fprintf)
	if err != nil {
		return
	}
}
