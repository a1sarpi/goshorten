package main

import (
	"github.com/a1sarpi/goshorten/api/handlers"
	"github.com/a1sarpi/goshorten/api/storage"
	"github.com/a1sarpi/goshorten/api/storage/memory"
	"log"
	"net/http"
)

func main() {
	var store storage.Storage
	store = memory.NewMemoryStorage()

	postHandler := handlers.NewPostHandler(store)
	getHandler := handlers.NewGetHandler(store)

	http.HandleFunc("/shorten", postHandler.HandleShorten)
	http.HandleFunc("/", getHandler.HandleRedirect)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
