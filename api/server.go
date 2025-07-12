package main

import (
	"log"
	"net/http"
	"strconv"
	"github.com/Rach17/Go-RSS-Aggregator/service"
	"github.com/Rach17/Go-RSS-Aggregator/utils"
)

type Server struct {
	Port        int
	Router      *http.ServeMux
	UserService *service.UserService
	AuthService *service.AuthService
	FeedService  *service.FeedService
	FeedPostService *service.FeedPostService
}

func NewServer(port int, userService *service.UserService, authService *service.AuthService, feedService *service.FeedService, feedPostService *service.FeedPostService) *Server {
	return &Server{
		Port:        port,
		Router:      http.NewServeMux(),
		UserService: userService,
		AuthService: authService,
		FeedService: feedService,
		FeedPostService: feedPostService,
	}

}

func (s *Server) Start() error {
	s.RegisterHandler()
	log.Printf("Starting server on port %d", s.Port)
	return http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Router)
}

func (s *Server) RegisterHandler() {
	s.Router.HandleFunc("GET /api/health", Chain(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}, corsMiddleware))
	AuthMiddleware := NewAuthMiddleware(s.AuthService)
	UserHandler := NewUserHandler(s.UserService)
	FeedHandler := NewFeedHandler(s.FeedService, s.UserService)
	FeedPostHandler := NewFeedPostHandler(s.FeedService, s.UserService, s .FeedPostService)

	s.Router.HandleFunc("POST /api/users", Chain(UserHandler.handleCreateUser, corsMiddleware))
	s.Router.HandleFunc("GET /api/users", Chain(UserHandler.handlerGetUserByAPIKey, AuthMiddleware.authMiddleware, corsMiddleware))

	s.Router.HandleFunc("POST /api/feed", Chain(FeedHandler.handleCreateFeed, AuthMiddleware.authMiddleware, corsMiddleware))
	s.Router.HandleFunc("GET /api/feed", Chain(FeedHandler.handleGetFeed, AuthMiddleware.authMiddleware, corsMiddleware))
	s.Router.HandleFunc("GET /api/feeds", Chain(FeedHandler.handleGetFeeds, AuthMiddleware.authMiddleware, corsMiddleware))
	s.Router.HandleFunc("POST /api/following", Chain(FeedHandler.handleFollowFeed, AuthMiddleware.authMiddleware, corsMiddleware))

	s.Router.HandleFunc("GET /api/feedposts", Chain(FeedPostHandler.handleGetFeedPost, AuthMiddleware.authMiddleware, corsMiddleware))
}