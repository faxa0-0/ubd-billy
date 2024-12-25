package middleware

import (
	"billy/utils"
	"billy/utils/jwt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func Logger(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("start", "method", r.Method, "path", r.URL.Path)
		defer slog.Info("end", "method", r.Method, "path", r.URL.Path)

		next(w, r)
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		token := parts[1]
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired access token")
			return
		}
		id, ok := (*claims)["sub"].(float64)
		if !ok {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}
		role, ok := (*claims)["role"].(string)
		if !ok {
			utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}
		r.Header.Set("X-User-ID", strconv.FormatFloat(id, 'f', -1, 64))
		r.Header.Set("X-User-ROLE", role)
		next(w, r)
	}
}
