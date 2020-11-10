package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
    log.Println("Starting server")
    buildHandler := http.FileServer(http.Dir("client/build"))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        buildHandler.ServeHTTP(w, r)
    })

    // Start the server, defaulting to port 3000.
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    log.Println("Running server on port", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
