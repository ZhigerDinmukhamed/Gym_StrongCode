package main

import (
	"log"
	"net/http"
	"os"
	
	"gym-api/router"
	"gym-api/store"
)

func main() {
	// ENV DEFAULT
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "CHANGE_ME_SECRET")
	}

	st, err := store.NewStore("gym.db")
	if err != nil {
		log.Fatalf("store init error: %v", err)
	}

	if err := st.InitSchema(); err != nil {
		log.Fatalf("schema init error: %v", err)
	}

	r := router.NewRouter(st)

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}
