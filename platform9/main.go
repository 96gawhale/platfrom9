package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Serve static files (HTML, CSS, JS)
    fs := http.FileServer(http.Dir("."))
    http.Handle("/", fs)

    // Handle requests for Calculator
    http.HandleFunc("/cal", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "http://localhost:8081", http.StatusFound)
    })

    // Handle requests for Todo List
    http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "http://localhost:8082", http.StatusFound)
    })

    // Start server on port 8080
    fmt.Println("Main Server started on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

