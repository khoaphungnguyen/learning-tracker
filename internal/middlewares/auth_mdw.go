package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khoaphungnguyen/learning-tracker/internal/auth"
)

// AuthMiddleware is a middleware that validates token and authorizes users
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
	
		// Get the Authorization header from the request
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			// If the Authorization header is not present, return a 403 status code
			c.JSON(403, "No Authorization header provided")
			c.Abort()
			return
		}
		// Split the Authorization header to get the token
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) == 2 {
			// Trim the token
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			// If the token is not in the correct format, return a 400 status code
			c.JSON(400, "Incorrect Format of Authorization Token")
			c.Abort()
			return
		}

		// Create a JwtWrapper with the secret key and issuer
		jwtWrapper := auth.JwtWrapper{
			SecretKey: secretKey,
			Issuer:    "AuthService",
		}

		// Validate the token
		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}
		// Set the claims in the context
		c.Set("email", claims.Email)
		c.Next()
	}
}
