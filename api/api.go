package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"

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

type GetResponse struct {
	Value string `json:"value"`
}

type DeleteRequest struct {
	Key string `json:"key"`
}

type AllResponse struct {
	Keys []string `json:"keys"`
}

type ServerOpts struct {
	ListenAddr string
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewApiServer(apiServerOpts ServerOpts, cache cache.Cacher) *Server {
	return &Server{
		ServerOpts: apiServerOpts,
		cache:      cache,
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/set", s.SetValue)
	router.HandleFunc("/get", s.GetValue)
	router.HandleFunc("/delete", s.DeleteValue)
	router.HandleFunc("/all", s.AllKeys)

	router.Use(loggingMiddleware)

	log.Printf("server starting on port [%s]\n", s.ListenAddr)
	log.Fatal(http.ListenAndServe(s.ListenAddr, router))
}

func (s *Server) WriteJson(w http.ResponseWriter, status int, v any) error {
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

func (s *Server) SetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = s.WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	req := new(SetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_ = s.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.cache.Set([]byte(req.Key), []byte(req.Value), req.Ttl); err != nil {
		_ = s.WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.WriteJson(w, http.StatusCreated, nil)
}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = s.WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	req := new(GetRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_ = s.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := new(GetResponse)
	value, err := s.cache.Get([]byte(req.Key))
	if err != nil {
		_ = s.WriteJson(w, http.StatusNotFound, err.Error())
		return
	}

	resp.Value = string(value)

	_ = s.WriteJson(w, http.StatusOK, resp)
}

func (s *Server) DeleteValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = s.WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	req := new(DeleteRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		_ = s.WriteJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.cache.Delete([]byte(req.Key)); err != nil {
		_ = s.WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = s.WriteJson(w, http.StatusOK, nil)
}

func (s *Server) AllKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = s.WriteJson(w, http.StatusMethodNotAllowed, nil)
		return
	}

	keysAsBytes, err := s.cache.All()
	if err != nil {
		_ = s.WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(AllResponse)
	keysAsStrings := make([]string, len(keysAsBytes))
	for i, key := range keysAsBytes {
		keysAsStrings[i] = string(key)
	}

	resp.Keys = keysAsStrings
	sort.Strings(resp.Keys)

	_ = s.WriteJson(w, http.StatusOK, resp)
}
