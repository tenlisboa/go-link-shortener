package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tenlisboa/go-link-shortener/internal/pkg/otel"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
	"github.com/tenlisboa/go-link-shortener/pkg/handlers"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ctx = context.Background()

func setupOTel() {
	serviceName := "go-link-shortener"
	serviceVersion := "0.1.0"
	otelShutdown, err := otel.SetupOTelSDK(ctx, serviceName, serviceVersion)
	if err != nil {
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	setupOTel()

	client := db.GetClient(ctx)
	r := mux.NewRouter()

	shortenHandler := handlers.NewShortenLinkHandler(handlers.ShortenLinkInput{
		Ctx: ctx,
		Dbc: *client,
	})
	sh := otelhttp.WithRouteTag("/short", http.HandlerFunc(shortenHandler.Retrieve))
	r.Handle("/short", sh).Methods("POST")

	retrieveHandler := handlers.NewRetrieveLinkHandler(handlers.RetrieveLinkInput{
		Ctx: ctx,
		Dbc: *client,
	})
	rh := otelhttp.WithRouteTag("/{hash}", http.HandlerFunc(retrieveHandler.Retrieve))
	r.Handle("/{hash}", rh)

	fmt.Println("Start to listen server at 8080")

	err = http.ListenAndServe(":8080", r)
	panic(err)
}
