package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const port = "7778"

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadFile("errors.json")
	fmt.Fprintln(w, string(body))
}

func main() {
	http.HandleFunc("/errors", handler)
	address := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Serving on address %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
