package handlers

import (
	"billy/models"
	"billy/storage"
	"billy/utils"
	"billy/utils/jwt"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.storage.GetUserByLogin(req.Login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, "Invalid credentials")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Pass))
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if user.Role != models.Admin {
		_, err = h.storage.EmulateUsageGathering(user.ID)
		if err != nil {
			slog.Warn("Failed to emulate usage gathering")
		}
	}

	token, err := jwt.CreateAccessToken(user.ID, user.Role)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Login failed, please try again")
		return
	}
	refresh_token, err := jwt.CreateRefreshToken(user.ID, user.Role)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Login failed, please try again")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
		Secure:   false, // need TLS `Let's Encrypt` maybe???
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	utils.SuccessResponse(w, http.StatusOK, map[string]string{"token": token}, "User fetched successfully")
}

func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "No refresh token found")
		return
	}

	refreshToken := refreshCookie.Value
	claims, err := jwt.ValidateToken(refreshToken)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired refresh token -validate")
		return
	}

	id, ok := (*claims)["sub"].(float64)
	if !ok {
		utils.ErrorResponse(w, http.StatusInternalServerError, "invalid or missing 'sub' claim in token")
		return
	}
	role, ok := (*claims)["role"].(string)
	if !ok {
		utils.ErrorResponse(w, http.StatusInternalServerError, "invalid or missing 'role' claim in token")
		return
	}

	accessToken, err := jwt.CreateAccessToken(int(id), models.Role(role))
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{"token": accessToken}, "Token refreshed successfully")
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("refresh_token")
	if err != nil {
		utils.SuccessResponse(w, http.StatusOK, nil, "User logged out successfully")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})

	utils.SuccessResponse(w, http.StatusOK, nil, "User logged out successfully")
}
