package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/OmprakashD20/refero-api/repository"
)

type APIServer struct {
	port string
	db   *repository.Queries
}

func NewAPIServer(port string, db *repository.Queries) *APIServer {
	return &APIServer{port, db}
}

func (s *APIServer) Run() error {
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()

	// TODO: Implement Custom Logger
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN] [%s] | %d | %s | %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.StatusCode,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	}))

	app.Use(gin.Recovery())

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "You hit the Refero API Server",
		})
	})

	api := app.Group("/api/v1")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "You hit the API v1 route of Refero",
			})
		})
	}

	log.Printf("Server is running on PORT %s", s.port)

	return app.Run(fmt.Sprintf(":%s", s.port))
}
