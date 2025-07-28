package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"shop/config"
	_ "shop/docs"
	"shop/internal/controller/http/handler"
	"shop/internal/usecase"
	"shop/pkg/logger"
	"shop/pkg/minio"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// NewRouter -.
// Swagger spec:
// @title       Ccenter News API
// @description This is a sample server Ccenter News server.
// @version     1.0
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, minio *minio.MinIO) {
	// Options
	engine.Use(gin.Logger())
	//engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, *minio)

	// Initialize Casbin enforcer

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend domenini yozish
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	engine.Use(TimeoutMiddleware(5 * time.Second))
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	engine.Use(cors.Default())
	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routes

	// auth
	// engine.POST("/auth/login", handlerV1.Login)

	engine.POST("/img-upload", handlerV1.UploadFile)

	users := engine.Group("/users")
	{
		users.POST("/login", handlerV1.Login)
		users.GET("/get", handlerV1.GetByIdUser)
		users.GET("/list", handlerV1.GetAllUsers)
		users.POST("/create", handlerV1.Register)
		users.PUT("/update", handlerV1.UpdateUser)
		users.DELETE("/delete", handlerV1.DeleteUser)
	}

	category := engine.Group("/category")
	{
		category.GET("/get", handlerV1.GetByIdCategory)
		category.GET("/list", handlerV1.GetAllCategories)
		category.POST("/create", handlerV1.CreateCategory)
		category.PUT("/update", handlerV1.UpdateCategory)
		category.DELETE("/delete", handlerV1.DeleteCategory)
	}

	shoes := engine.Group("/shoes")
	{
		shoes.GET("/get", handlerV1.GetByIdShoes)
		shoes.GET("/list", handlerV1.GetAllShoes)
		shoes.POST("/create", handlerV1.CreateShoes)
		shoes.PUT("/update", handlerV1.UpdateShoes)
		shoes.DELETE("/delete", handlerV1.DeleteShoes)
	}
}
