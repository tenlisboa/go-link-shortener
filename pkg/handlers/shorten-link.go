package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tenlisboa/go-link-shortener/internal/pkg/jsonp"
	"go.opentelemetry.io/otel"
	"net/http"
	"os"

	"github.com/tenlisboa/go-link-shortener/internal/pkg/random"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
	"gopkg.in/go-playground/validator.v9"
)

var (
	tracer = otel.Tracer("shorten")
)

type ShortenLinkHandler struct {
	ctx context.Context
	dbc db.Client
}

type ShortenLinkInput struct {
	Dbc db.Client
	Ctx context.Context
}

type shortenLinkEntity struct {
	URL string `json:"url" validate:"required,uri"`
}

func NewShortenLinkHandler(input ShortenLinkInput) *RetrieveLinkHandler {
	return &RetrieveLinkHandler{
		ctx: input.Ctx,
		dbc: input.Dbc,
	}
}

func (sl *RetrieveLinkHandler) Store(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "shorten")
	defer span.End()

	body := jsonp.ToStruct[shortenLinkEntity](r.Body)

	v := validator.New()
	err := v.Struct(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	hash := random.Hash(6)
	err = sl.dbc.Store(hash, body.URL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%v/%v", os.Getenv("APP_URL"), hash)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(shortenLinkEntity{
		URL: url,
	})
}
