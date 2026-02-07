package integration

import (
	"context"
	"testing"
	"time"
	"github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()
	poolCtx, cancel := context.WithTimeout(ctx, 15 * time.Second)
	defer cancel()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16.0",
		postgres.WithDatabase("testDB"),
		postgres.WithPassword("testDB"),
		postgres.WithUsername("testDB"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5 * time.Second)),
	)

	if err != nil {
		t.Fatalf("could not start postgres container: %s", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("could not get connection string: %s", err)
	}

	pool, err := pgxpool.New(poolCtx, connStr)
	if err != nil {
		t.Fatalf("could not create pool: %s", err)
	}

	createSchema(t, pool)

	cleanup := func() {
		pool.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("could not terminate container: %s", err)
		}
	}

	return pool, cleanup
}

func createSchema(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()

	schema := `CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL
	);`

	_, err := pool.Exec(ctx, schema)
	if err != nil {
		t.Fatalf("could not create schema: %s", err)
	}
}