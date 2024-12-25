package handlers

import (
	"billy/models"
	"billy/storage"
	"billy/utils"
	"errors"
	"net/http"
	"strconv"
)

func (h *Handler) GetUsageHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-ROLE")
	if role == "" {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Role not found in request")
		return
	}
	tokenUserID := r.Header.Get("X-User-ID")
	queryUserID := r.URL.Query().Get("id")

	if models.Role(role) != models.Admin {
		if tokenUserID != queryUserID {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Not your level request")
			return
		}
	}
	id, err := strconv.Atoi(queryUserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	usage, err := h.storage.GetUsageByUserID(id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, "Usage not found")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch usage")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, usage, "Usage found")
}
