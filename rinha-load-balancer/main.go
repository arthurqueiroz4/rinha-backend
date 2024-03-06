package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type LoadBalancer interface {
	GetBackend() url.URL
}

type RoundRobin struct {
	Backends []url.URL
	mux      sync.Mutex
	counter  int
}

func NewRoundRobin(backends []url.URL) *RoundRobin {
	return &RoundRobin{
		Backends: backends,
	}
}

func (s *RoundRobin) GetBackend() url.URL {
	s.mux.Lock()
	s.counter = (s.counter + 1) % len(s.Backends)
	defer s.mux.Unlock()
	return s.Backends[s.counter]
}

type Handler struct {
	LB LoadBalancer
}

func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	backend := h.LB.GetBackend()
	backend.Path = r.URL.Path

	req, err := http.NewRequest(r.Method, backend.String(), r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func main() {
	api01 := "localhost:" + os.Getenv("API01_PORT")
	api02 := "localhost:" + os.Getenv("API02_PORT")

	rr := NewRoundRobin([]url.URL{
		{Scheme: "http", Host: api01},
		{Scheme: "http", Host: api02},
	})
	h := &Handler{LB: rr}

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(h.handleRequest)

	log.Print("LOAD BALANCER RODANDO NA PORTA 9999 /// API01: " + api01 + " -- API02: " + api02)
	log.Fatal(http.ListenAndServe(":9999", r))
}
