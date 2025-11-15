package main

import (
	"log"
	"net/http"

	_ "github.com/ZhigerDinmukhamed/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	store, err := NewStore("gym.db")
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}
	defer store.DB.Close()

	if err := store.InitSchema(); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	r := NewRouter(store)

	// Добавляем Swagger UI
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
