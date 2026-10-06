package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

func VerifyJWT(next httprouter.Handle) httprouter.Handle {
	_ = godotenv.Load()

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		jwtKey = "sementara123!"
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenString := ""
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				tokenString = authHeader
			}
		} else if cookie, err := r.Cookie("authToken"); err == nil && cookie.Value != "" {
			tokenString = cookie.Value
		}

		if tokenString == "" {
			helper.WriteJSONError(w, http.StatusUnauthorized, "missing authorization token")
			return
		}

		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			helper.WriteJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		r.Header.Set("User-Name", claims.Username)
		next(w, r, ps)
	}
}
