package transport

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

// SetupRoutes sets up all the routes for the API
func (h *NetHandler) SetupRoutes(router *http.ServeMux) {
	// Routes for user operations
	router.HandleFunc("/users", h.authMiddleware(h.handleUsers))
	router.HandleFunc("/users/add", h.handleAddUsers)

	// Routes for auth operation
	router.HandleFunc("/auth/signin", h.handleSignIn)

	// Routes for goals operations
	router.HandleFunc("/goals", h.authMiddleware(h.handleGoals))
	router.HandleFunc("/goals/add", h.authMiddleware(h.handleNewGoal))

	// Routes for entry operations
	router.HandleFunc("/entries", h.authMiddleware(h.handleEntries))
	router.HandleFunc("/entries/add", h.authMiddleware(h.handleNewEntry))

	// Routes for file operations
	router.HandleFunc("/files", h.authMiddleware(h.handleFiles))
	router.HandleFunc("/files/add", h.authMiddleware(h.handleNewFile))

}

func (h *NetHandler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			http.Error(w, "Missing or malformed JWT", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(authHeader[1], &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(h.JWTKey), nil
		})

		if err != nil {
			http.Error(w, "Invalid or expired JWT: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid or expired JWT", http.StatusUnauthorized)
			return
		}

		userID := claims.Audience

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
