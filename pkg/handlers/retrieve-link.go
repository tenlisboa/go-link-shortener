package handlers

import (
	"context"
	localotel "github.com/tenlisboa/go-link-shortener/internal/pkg/otel"
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
	tr := localotel.TraceIt(sl.ctx, "retrieve")
	ctx, span := tr.Start(sl.ctx, "retrieve")
	defer span.End()

	params := mux.Vars(r)
	hash := params["hash"]

	_, dbspan := tr.Start(ctx, "get from db")
	url, err := sl.dbc.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	dbspan.End()

	http.Redirect(w, r, url, http.StatusSeeOther)
}
