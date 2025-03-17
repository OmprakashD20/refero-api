package middlewares

import (
	"log"
	"net/http"

	"github.com/OmprakashD20/refero-api/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors which occurs during request processing.
// If an error is found in the context, it logs the error and responds with an appropriate HTTP status code and message.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			log.Println(err.Err.Error())

			if e, ok := err.Err.(*errors.HTTPError); ok {
				c.JSON(e.StatusCode, gin.H{
					"error": e.ErrorMsg,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}

			c.Abort()
		}
	}
}
