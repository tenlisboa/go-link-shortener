package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tenlisboa/go-link-shortener/internal/pkg/jsonp"
	"net/http"

	"github.com/tenlisboa/go-link-shortener/internal/pkg/random"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
	"gopkg.in/go-playground/validator.v9"
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

	url := fmt.Sprintf("%v/%v", "http://localhost:8080", hash)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenLinkEntity{
		URL: url,
	})
}
