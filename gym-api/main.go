package main

import (
	"log"
	"net/http"
)

func main() {
	// initialize DB/store
	store, err := NewStore("gym.db")
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}
	defer store.DB.Close()

	// ensure schema
	if err := store.InitSchema(); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	// create router and handlers
	r := NewRouter(store)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
