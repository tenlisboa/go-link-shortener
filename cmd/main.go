package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
	"github.com/tenlisboa/go-link-shortener/pkg/handlers"
)

var ctx = context.Background()

func main() {

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
