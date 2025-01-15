package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"redirector/models"
)

// TokenAuthMiddleware validates Bearer tokens using a KV store and blocks unauthorized requests.
func TokenAuthMiddleware(kv models.KV) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Grab the "Authorization" header and split it at the first space.
		authData := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		// A valid header will have a length of 2.
		if len(authData) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// A valid header starts with "Bearer"
		if authData[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Try to get the requested auth token from the database, if it exists the token is valid.
		data, err := kv.Get([]byte(authData[1]))
		if err != nil || data == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
