package main

import (
	"dbaas/handler"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", handler.GetHandler)
	http.ListenAndServe(":8080", nil)
}
