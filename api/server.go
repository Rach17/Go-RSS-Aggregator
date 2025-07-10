package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Rach17/Go-RSS-Aggregator/service"
)

type Server struct {
	Port        int
	Router      *http.ServeMux
	UserService *service.UserService
	AuthService *service.AuthService
	FeedService  *service.FeedService
}

func NewServer(port int, userService *service.UserService, authService *service.AuthService, feedService *service.FeedService) *Server {
	return &Server{
		Port:        port,
		Router:      http.NewServeMux(),
		UserService: userService,
		AuthService: authService,
		FeedService: feedService,
		
	}

}

func (s *Server) Start() error {
	s.RegisterHandler()
	log.Printf("Starting server on port %d", s.Port)
	return http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Router)
}

func (s *Server) RegisterHandler() {
	s.Router.HandleFunc("GET /api/health", Chain(handlerReadiness, corsMiddleware))
	AuthMiddleware := NewAuthMiddleware(s.AuthService)
	UserHandler := NewUserHandler(s.UserService)

	s.Router.HandleFunc("POST /api/users", Chain(UserHandler.handleCreateUser, corsMiddleware))
	s.Router.HandleFunc("GET /api/users", Chain(UserHandler.handlerGetUserByAPIKey, AuthMiddleware.authMiddleware, corsMiddleware))

	s.Router.HandleFunc("POST /api/feed", Chain(NewFeedHandler(s.FeedService).handleCreateFeed, AuthMiddleware.authMiddleware, corsMiddleware))
	s.Router.HandleFunc("GET /api/feed", Chain(NewFeedHandler(s.FeedService).handleGetFeed, AuthMiddleware.authMiddleware, corsMiddleware))
}
