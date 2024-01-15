package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/preseed.cfg", handleVMFileRequest)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleVMFileRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, World")
}
