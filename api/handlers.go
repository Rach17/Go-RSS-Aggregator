package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rach17/Go-RSS-Aggregator/service"
	"github.com/Rach17/Go-RSS-Aggregator/utils"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (handler *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var params parameters = parameters{}
	var decoder *json.Decoder = json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err := handler.UserService.CreateUser(r.Context(), params.Username, params.Password)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create user: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (handler *UserHandler) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	const userContextKey contextKey = "user"
	user := r.Context().Value(userContextKey)
	if user == nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, user)
}

type FeedHandler struct {
	FeedService *service.FeedService
}

func NewFeedHandler(feedService *service.FeedService) *FeedHandler {
	return &FeedHandler{
		FeedService: feedService,
	}
}

func (h *FeedHandler) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		URL string `json:"url"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create feed (business logic)
	feed, err := h.FeedService.CreateFeed(r.Context(), params.URL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create feed: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, feed)
}

func (h *FeedHandler) handleGetFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		URL string `json:"url"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feed, err := h.FeedService.GetFeedByURL(r.Context(), params.URL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get feeds: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, feed)
}
