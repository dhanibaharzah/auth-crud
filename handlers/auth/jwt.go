// auth/jwt.go
package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Define the secret key used to sign and verify JWT tokens.
var jwtKey = []byte("ramadhani-aa-bb-cc")

func GetCurrentUserIDFromContext(r *http.Request) int {
	// Extract the JWT token from the request's Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0 // No token found
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method (HS256 in this case)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0 // Invalid token
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user ID from the claims
		if userID, ok := claims["userID"].(int); ok {
			return userID
		}
	}

	return 0 // Invalid or missing user ID in the token
}
