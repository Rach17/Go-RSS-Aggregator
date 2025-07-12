package main 

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/Rach17/Go-RSS-Aggregator/service"
	"github.com/Rach17/Go-RSS-Aggregator/utils"
)

type FeedPostHandler struct {
	FeedService  *service.FeedService
	UserService  *service.UserService
	FeedPostService *service.FeedPostService
}

func NewFeedPostHandler(feedService *service.FeedService, userService *service.UserService, feedPostService *service.FeedPostService) *FeedPostHandler {
	return &FeedPostHandler{
		FeedService:  feedService,
		UserService:  userService,
		FeedPostService: feedPostService,
	}
}

func (h *FeedPostHandler) handleGetFeedPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		URL string `json:"url"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := r.Context().Value(userContextKey)
	if user == nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get feed posts (business logic)
	posts, err := h.FeedPostService.GetFeedPosts(r.Context(), params.URL)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get feed posts: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, posts)
}