package controller

import (
	model "gin-demo/internal/models"
	service "gin-demo/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	service *service.NewsService
}

func NewNewsController(service *service.NewsService) *NewsController {
	return &NewsController{service}
}

func (c *NewsController) GetAll(ctx *gin.Context) {
	data, _ := c.service.GetAllNews()
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func (c *NewsController) Create(ctx *gin.Context) {
	var news model.News
	if err := ctx.ShouldBindJSON(&news); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.service.CreateNews(&news)
	ctx.JSON(200, gin.H{"message": "created"})
}
