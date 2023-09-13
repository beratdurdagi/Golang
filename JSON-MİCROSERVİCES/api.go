package main

import (
	"context"
	"math/rand"

	"encoding/json"
	"net/http"

	"github.com/karalakrepp/Golang/BasicMicroservices/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func makeAPIFunc(fnc APIFunc) http.HandlerFunc {

	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx = context.WithValue(ctx, "requestID", rand.Intn(100000000))
		if err := fnc(context.Background(), w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}

	}
}

type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJSONAPIServer(svc PriceFetcher, listenAddr string) *JSONAPIServer {

	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {

	http.HandleFunc("/", makeAPIFunc(s.handlePrice))
	http.ListenAndServe(s.listenAddr, nil)
}

func (s *JSONAPIServer) handlePrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ticker := r.URL.Query().Get("ticker")

	price, err := s.svc.FetchPrice(ctx, ticker)

	if err != nil {
		return err
	}

	resp := types.ClientResponse{
		Ticker: ticker,
		Price:  price,
	}
	return writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
