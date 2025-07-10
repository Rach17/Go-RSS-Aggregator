package api

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/Rach17/Go-RSS-Aggregator/utils"
	"github.com/Rach17/Go-RSS-Aggregator/service"
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