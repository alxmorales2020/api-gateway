package main

import (
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

    log.Println("Gateway listening on :8080")
    http.ListenAndServe(":8080", r)
}
