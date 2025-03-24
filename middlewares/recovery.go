package middlewares

import (
	"log"
	"net/http"

	// "github.com/OmprakashD20/refero-api/errors"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from panics that occur during request processing.
// If a panic occurs, it sets an internal server error in the context.
// This prevents the server from getting crashed due to a panic.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC RECOVERED: %v\n", rec)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()

		c.Next()
	}
}
