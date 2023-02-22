package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/TheLazyLemur/cacheengine/cache"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Ttl   int64  `json:"ttl"`
}

type GetRequest struct {
	Key string `json:"key"`
}

type GetResonse struct {
	Value string `json:"value"`
}

type DeleteRequest struct {
	Key string `json:"key"`
}

type ApiServerOpts struct {
	ListenAddr string
}

type ApiServer struct {
	ApiServerOpts
	cacher cache.Cacher
}

func NewApiServer(apiServerOpts ApiServerOpts, cache cache.Cacher) *ApiServer {
	return &ApiServer{
		ApiServerOpts: apiServerOpts,
		cacher:        cache,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/set", s.SetValue)
	router.HandleFunc("/get", s.GetValue)
	router.HandleFunc("/delete", s.DeleteValue)

	router.Use(loggingMiddleware)

	log.Printf("server starting on port [%s]\n", s.ListenAddr)
	log.Fatal(http.ListenAndServe(s.ListenAddr, router))
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ContextInformation struct {
	CtxKey string
	Uuid   string
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxInf := new(ContextInformation)
		ctxInf.Uuid = uuid.New().String()

		ctxInf.CtxKey = "requestInf"
		ctx := context.WithValue(r.Context(), ctxInf.CtxKey, ctxInf)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *ApiServer) SetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	req := new(SetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_ = WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.cacher.Set([]byte(req.Key), []byte(req.Value), req.Ttl); err != nil {
		_ = WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJson(w, http.StatusCreated, nil)
}

func (s *ApiServer) GetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	req := new(GetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_ = WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := new(GetResonse)
	value, err := s.cacher.Get([]byte(req.Key))
	if err != nil {
		_ = WriteJson(w, http.StatusNotFound, err.Error())
		return
	}

	resp.Value = string(value)

	_ = WriteJson(w, http.StatusOK, resp)
}

func (s *ApiServer) DeleteValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	req := new(DeleteRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_ = WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.cacher.Delete([]byte(req.Key)); err != nil {
		_ = WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJson(w, http.StatusOK, nil)
}
