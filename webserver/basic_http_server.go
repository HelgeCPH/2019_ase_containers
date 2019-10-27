package main

import (
    "fmt"
    "log"
    "net/http"
)

func hejVerdenHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hej verden!\n")
}

func main() {
    port := 8080

    http.HandleFunc("/", hejVerdenHandler)

    log.Printf("Server starting on port %v\n", port)
    http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
