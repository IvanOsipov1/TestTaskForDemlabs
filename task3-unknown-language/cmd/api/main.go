package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	apihttp "github.com/yourname/inventory-go/internal/http"
	"github.com/yourname/inventory-go/internal/store"
)

func main() {
	// SQLite файл в корне: inventory.db
	// DSN с включённым foreign_keys
	dsn := "file:inventory.db?cache=shared&_pragma=foreign_keys(1)"
	if v := os.Getenv("DB_DSN"); v != "" {
		dsn = v
	}

	db, err := store.Open(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	st := store.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      apihttp.Router(st),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
