package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OmprakashD20/refero-api/database"
	"github.com/OmprakashD20/refero-api/middlewares"
	"github.com/OmprakashD20/refero-api/services/category"
	"github.com/OmprakashD20/refero-api/services/links"
)

type APIServer struct {
	port string
	conn *pgxpool.Pool
}

func NewAPIServer(port string, conn *pgxpool.Pool) *APIServer {
	return &APIServer{port, conn}
}

func (s *APIServer) Run() error {
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()

	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN] [%s] | [%s] %d | %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.StatusCode,
			param.ClientIP,
			param.Path,
		)
	}))

	app.Use(middlewares.RecoveryMiddleware())
	app.Use(middlewares.ErrorHandler())

	// Handle 404 API routes
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "route not found",
		})
	})

	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "you hit the Refero API Server",
		})
	})

	api := app.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "you hit the v1 API route of Refero",
			})
		})

		txnStore := database.NewTransactionStore(s.conn)

		// Category Routes
		categoryStore := category.NewStore(s.conn)
		categoryService := category.NewService(categoryStore)
		categoryService.SetupCategoryRoutes(api.Group("/category"))

		// Link Routes
		linkStore := links.NewStore(s.conn)
		LinkService := links.NewService(linkStore, txnStore)
		LinkService.SetupLinkRoutes(api.Group("/link"))
	}

	log.Printf("Server is running on PORT %s", s.port)

	return app.Run(fmt.Sprintf(":%s", s.port))
}
