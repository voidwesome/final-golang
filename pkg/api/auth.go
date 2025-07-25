package api

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret") // will be replaced by password-derived value in runtime

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			// no auth required
			next(w, r)
			return
		}
		// derive secret from password
		jwtSecret = []byte(pass) // simple & enough for the task

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*Claims)
		if !ok || claims.Hash != pass {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}

func makeToken(pass string) (string, error) {
	jwtSecret = []byte(pass)
	claims := &Claims{
		Hash: pass,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}
