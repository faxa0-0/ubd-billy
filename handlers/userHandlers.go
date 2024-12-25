package handlers

import (
	"billy/models"
	"billy/storage"
	"billy/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-ROLE")
	if models.Role(role) != models.Admin {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Not your level request")
		return
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if user.Name == "" || user.Login == "" || user.Pass == "" || user.Plan == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	if user.Role == "" {
		user.Role = models.Customer
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	user.Pass = string(hash)

	id, err := h.storage.CreateUser(user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, map[string]int{"id": id}, "User created successfully")
}

func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-ROLE")
	if models.Role(role) != models.Admin {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Not your level request")
		return
	}
	users, err := h.storage.GetUsers()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	if len(users) == 0 {
		utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{"users": []interface{}{}}, "No users found")
		return
	}

	for i := range users {
		users[i].Pass = ""
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{"users": users}, "Users fetched successfully")
}

func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-ROLE")
	if role == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Role not found in request")
		return
	}

	tokenUserID := r.Header.Get("X-User-ID")
	token_sub, err := strconv.Atoi(tokenUserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if models.Role(role) != models.Admin {
		if token_sub != id {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Not your level request")
			return
		}
	}

	user, err := h.storage.GetUserById(id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, "User not found")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}
	user.Pass = ""

	utils.SuccessResponse(w, http.StatusOK, user, "User fetched successfully")
}
