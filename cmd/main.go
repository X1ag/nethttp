package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"nethttppractice/internal/api/handlers"
	"nethttppractice/internal/item"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DB_DSN") 

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	poolRepo := itemrepo.NewPgRepo(pool)

	handlers := handlers.NewItemHandler(poolRepo)

	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	slog.Default().Info("Connected to database")

	serv := http.ServeMux{}

	serv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	serv.HandleFunc("GET /items", handlers.GetItems) 
	serv.HandleFunc("POST /items", handlers.InsertItems)

	slog.Default().Info("Server is running")
	http.ListenAndServe(":8080", &serv)
}