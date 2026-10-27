package main

import (
	"context"
	"log"
	"net/http"

	"movank/internal/config"
	"movank/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("migraciones: %v", err)
	}
	log.Println("migraciones OK")

	// La capa httpapi añadirá las rutas en el siguiente paso.
	mux := http.NewServeMux()

	log.Printf("escuchando en :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("servidor: %v", err)
	}
}