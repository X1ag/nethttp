package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"nethttppractice/internal/usecases"
	postgres "nethttppractice/internal/repository"

	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN") 

	ctx := context.Background()
	poolCtx, cancel := context.WithTimeout(ctx, 10 * time.Second)

	defer cancel()

	pool, err := pgxpool.New(poolCtx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	poolRepo := postgres.NewPgRepo(pool)
	handlers := usecases.NewItemHandler(poolRepo)

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	slog.Default().Info("Connected to database")

	serv := http.ServeMux{}

	serv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	serv.HandleFunc("GET /items", handlers.GetItems) 
	serv.HandleFunc("POST /items", handlers.InsertItem)
	serv.HandleFunc("DELETE /items/{id}", handlers.DeleteItem)

	slog.Default().Info("Server is running")

	http.ListenAndServe(":8080", &serv)
}