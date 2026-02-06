package router

import (
	controller "gin-demo/internal/controllers"
	"gin-demo/internal/middleware"
	repository "gin-demo/internal/repositories"
	service "gin-demo/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	// 有 Global Error Handler 時不要用 gin.Default()，因為要掌控 error flow
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())

	newsRepo := repository.NewNewsRepository(db)
	logRepo := repository.NewNewsLogRepository(db)
	newsService := service.NewNewsService(db, newsRepo, logRepo)
	newsController := controller.NewNewsController(newsService)

	api := r.Group("/api")

	news := api.Group("/news")
	{
		news.GET("", newsController.GetAll)
		news.GET(":id", newsController.GetByID)
		news.POST("/batch", middleware.JWTAuth(), newsController.BatchImport)
		news.POST("", middleware.JWTAuth(), newsController.Create)
		news.PUT(":id", middleware.JWTAuth(), newsController.Update)
		news.DELETE(":id", middleware.JWTAuth(), newsController.Delete)
	}

	return r
}
