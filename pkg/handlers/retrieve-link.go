package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
)

type RetrieveLinkHandler struct {
	ctx context.Context
	dbc db.Client
}

type RetrieveLinkInput struct {
	Dbc db.Client
	Ctx context.Context
}

func NewRetrieveLinkHandler(input RetrieveLinkInput) *RetrieveLinkHandler {
	return &RetrieveLinkHandler{
		ctx: input.Ctx,
		dbc: input.Dbc,
	}
}

func (sl *RetrieveLinkHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]

	url, err := sl.dbc.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
