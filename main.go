package main

import (
	"fmt"
	"net/http"
)

// Helloworld handler function that will be called when someone visits /api/hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! The Lemonform backend is alive.")
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	fmt.Println("Starting server on http://localhost:8080 ...")

	// Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}