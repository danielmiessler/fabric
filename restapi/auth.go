package restapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerApiKey := c.GetHeader("X-API-Key")

		if headerApiKey != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Wrong or missing API Key")})
			return
		}

		c.Next()
	}
}
