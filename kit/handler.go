package kit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AbortWithStatusJSON is a helper function that calls `Abort()` and then `JSON` internally.
// This method stops the chain, writes the status code and return a JSON body with HTTP status code and error message.
// It also sets the Content-Type as "application/json".
func AbortWithStatusJSON(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// NotFoundHandler is a helper function that calls server.AbortWithStatusJSON.
func NotFoundHandler(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
