package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zgiber/ports/internal/parser"
	"github.com/zgiber/ports/service"
)

type Service interface {
	UpdatePorts(ctx context.Context, feed service.Feed) error
	ListPorts(ctx context.Context, filter service.PortsFilter) ([]*service.Port, error)
}

func NewServer(ctx context.Context, s Service, addr string) *Server {
	server := &Server{}
	server.service = s
	server.m = http.NewServeMux()
	server.Server = http.Server{
		Addr:    addr,
		Handler: server.m,
	}

	server.initRoutes()
	return server
}

type Server struct {
	http.Server
	m       *http.ServeMux
	service Service
}

func (s *Server) initRoutes() {
	s.m.HandleFunc("/ports", s.handlePorts)
}

func (s *Server) handlePorts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getPortList(w, r)
	case http.MethodPost:
		s.postUpdateBatch(w, r)
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) postUpdateBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feed := parser.New(ctx, r.Body)
	if err := s.service.UpdatePorts(ctx, feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) getPortList(w http.ResponseWriter, r *http.Request) {
	maxItems, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	filter := service.PortsFilter{
		AfterID:  r.URL.Query().Get("from_id"),
		MaxItems: maxItems,
	}

	ports, err := s.service.ListPorts(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(ports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
