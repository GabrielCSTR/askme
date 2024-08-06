package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"

	"fmt"
	"os"

	"github.com/GabrielCSTR/askme.git/internal/api"
	"github.com/GabrielCSTR/askme.git/internal/store/pgstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Erro ao carregar variáveis de ambiente: %v", err)
	}

	ctx := context.Background() // Declare and initialize ctx variable

	pool, err := pgxpool.New(ctx, fmt.Sprintf(("host=%s port=%s user=%s password=%s dbname=%s"),
		os.Getenv("WSRS_DATABASE_HOST"),
		os.Getenv("WSRS_DATABASE_PORT"),
		os.Getenv("WSRS_DATABASE_USER"),
		os.Getenv("WSRS_DATABASE_PASSWORD"),
		os.Getenv("WSRS_DATABASE_NAME"),
	))

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Connected to database ⚡")

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	fmt.Println("Server running on port 8080 🚀")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
