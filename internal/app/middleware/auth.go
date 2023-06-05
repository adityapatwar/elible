// internal/app/middleware/auth.go
package middleware

// import (
// 	"elible/internal/app/utils"
// 	"elible/internal/config"
// 	"net/http"
// 	"strings"
// )

// func AuthMiddleware(next http.Handler, cfg *config.Config) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check if the Authorization header is set
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Authorization header required", http.StatusUnauthorized)
// 			return
// 		}

// 		// Check if the Authorization header is well-formed
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
// 			http.Error(w, "Invalid Authorization header format. Format is 'Bearer {token}'", http.StatusUnauthorized)
// 			return
// 		}

// 		// Validate the token
// 		token := tokenParts[1]
// 		claims, err := utils.ValidateJWTWithLocalSecret(token,cfg)
// 		if err != nil {
// 			http.Error(w, "Invalid Authorization token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Continue if the token is valid
// 		next.ServeHTTP(w, r)
// 	})
// }
