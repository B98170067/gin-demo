package router

import (
	controller "gin-demo/internal/controllers"
	repository "gin-demo/internal/repositories"
	service "gin-demo/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	repo := repository.NewNewsRepository(db)
	service := service.NewNewsService(repo)
	controller := controller.NewNewsController(service)

	api := r.Group("/api")

	news := api.Group("/news")
	{
		news.GET("", controller.GetAll)
		news.GET(":id", controller.GetByID)
		news.POST("", controller.Create)
		news.PUT(":id", controller.Update)
		news.DELETE(":id", controller.Delete)
	}

	return r
}
