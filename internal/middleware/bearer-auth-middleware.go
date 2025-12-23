package middleware

import (
	"net/http"
	"strings"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/gin-gonic/gin"
)

func BearerAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		token := strings.Split(bearerToken, " ")[1]
		tokenDto, err := frameworkutils.ParseJWT(token, []byte(jwtSecret))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set(frameworkconstants.TokenKey, tokenDto)

		c.Next()
	}
}
