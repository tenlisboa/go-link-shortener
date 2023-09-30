package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tenlisboa/go-link-shortener/internal/pkg/jsonp"
	"go.opentelemetry.io/otel"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os"

	localotel "github.com/tenlisboa/go-link-shortener/internal/pkg/otel"
	"github.com/tenlisboa/go-link-shortener/internal/pkg/random"
	"github.com/tenlisboa/go-link-shortener/pkg/db"
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

func NewShortenLinkHandler(input ShortenLinkInput) *ShortenLinkHandler {
	return &ShortenLinkHandler{
		ctx: input.Ctx,
		dbc: input.Dbc,
	}
}

func (sl *ShortenLinkHandler) Store(w http.ResponseWriter, r *http.Request) {
	tr := localotel.TraceIt(sl.ctx, "store")
	ctx, span := tr.Start(sl.ctx, "short")
	defer span.End()

	body := jsonp.ToStruct[shortenLinkEntity](r.Body)

	_, vspan := tr.Start(ctx, "validate")
	v := validator.New()
	err := v.Struct(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	vspan.End()

	_, hspan := tr.Start(ctx, "hash")
	hash := random.Hash(6)
	hspan.End()

	_, dbspan := tr.Start(ctx, "save on db")
	err = sl.dbc.Store(hash, body.URL)
	dbspan.End()

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
