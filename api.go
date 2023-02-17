package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TheLazyLemur/cacheengine/cache"
	"github.com/gorilla/mux"
)

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Ttl   int64  `json:"ttl"`
}

type GetRequest struct {
	Key   string `json:"key"`
}

type GetResonse struct {
	Value   string `json:"value"`
}

type DeleteRequest struct {
	Key   string `json:"key"`
}

type ApiServerOpts struct {
	ListenAddr string
}

type ApiServer struct {
	ApiServerOpts
	cacher      cache.Cacher
}

func NewApiServer(apiServerOpts ApiServerOpts, cache cache.Cacher) *ApiServer {
	return &ApiServer{
		ApiServerOpts: apiServerOpts,
		cacher:      cache,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/set", s.SetValue)
	router.HandleFunc("/get", s.GetValue)
	router.HandleFunc("/delete", s.DeleteValue)

	log.Printf("server starting on port [%s]\n", s.ListenAddr)
	log.Fatal(http.ListenAndServe(s.ListenAddr, router))
}

func (s *ApiServer) SetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := new(SetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.cacher.Set([]byte(req.Key), []byte(req.Value), req.Ttl); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ApiServer) GetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := new(GetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := new(GetResonse)
	value, err := s.cacher.Get([]byte(req.Key))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp.Value = string(value)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ =  json.NewEncoder(w).Encode(resp)
}

func (s *ApiServer) DeleteValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := new(DeleteRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.cacher.Delete([]byte(req.Key)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
