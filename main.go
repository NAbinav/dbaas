package main

import (
	"dbaas/db"
	"dbaas/handler"
	"net/http"
)

func main() {

	db.Init_DB()
	http.HandleFunc("GET /", handler.GetHandler)
	http.HandleFunc("POST /", handler.PostHandler)
	http.ListenAndServe(":8080", nil)
}
