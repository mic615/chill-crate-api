package auth

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenClaims struct {
	Sub       string `json:"sub"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Username  string `json:"preferred_username"`
}

func (a *Authenticator) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// truly missing
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing authorization header"})
			return
		}
		token, present := strings.CutPrefix(authHeader, "Bearer ")
		if !present || token == "" {
			// present but wrong format, or "Bearer " with no token
			c.AbortWithStatusJSON(401, gin.H{"error": "Malformed authorization header"})
			return
		}
		idToken, err := a.kc.OIDC.Verify(c.Request.Context(), token)
		if err != nil {
			log.Printf("token verification failed: %v", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		var claims TokenClaims
		if err = idToken.Claims(&claims); err != nil {
			log.Printf("failed to parse token claims: %v", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Failed to parse token claims"})
			return
		}
		user, err := a.resolveOrCreateUser(&claims)
		if err != nil {
			log.Printf("failed to resolve or create user: %v", err)
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to resolve or create user"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
