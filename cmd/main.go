package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tenlisboa/go-link-shortener/internal/pkg/otel"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
	"github.com/tenlisboa/go-link-shortener/pkg/handlers"
	"net/http"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	cleanup := otel.SetupOTelTracer(ctx)
	defer cleanup(ctx)

	client := db.GetClient(ctx)

	r := mux.NewRouter()

	shortenHandler := handlers.NewShortenLinkHandler(handlers.ShortenLinkInput{
		Ctx: ctx,
		Dbc: *client,
	})
	r.HandleFunc("/short", shortenHandler.Store).Methods("POST")

	retrieveHandler := handlers.NewRetrieveLinkHandler(handlers.RetrieveLinkInput{
		Ctx: ctx,
		Dbc: *client,
	})
	r.HandleFunc("/{hash}", retrieveHandler.Retrieve)

	fmt.Println("Start to listen server at 8080")

	http.ListenAndServe(":8080", r)
}
