package jwt

import (
	"billy/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

func CreateToken(id int, role models.Role, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  id,
			"role": role,
			"exp":  time.Now().Add(expiry).Unix(),
			"iat":  time.Now().Unix(),
		})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func CreateAccessToken(id int, role models.Role) (string, error) {
	return CreateToken(id, role, AccessTokenExpiration)
}

func CreateRefreshToken(id int, role models.Role) (string, error) {
	return CreateToken(id, role, RefreshTokenExpiration)
}

func ValidateToken(token string) (*jwt.MapClaims, error) {
	jwToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !jwToken.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := jwToken.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	expirationTime := int64((*claims)["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		return nil, errors.New("refresh token expired")
	}
	return claims, nil
}
