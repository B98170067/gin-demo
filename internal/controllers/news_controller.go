package controller

import (
	model "gin-demo/internal/models"
	service "gin-demo/internal/services"
	"net/http"
	"strconv"

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

func (c *NewsController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(200, gin.H{"data": data})
}

func (c *NewsController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var news model.News
	if err := ctx.ShouldBindJSON(&news); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	news.ID = uint(id)
	c.service.Update(&news)
	ctx.JSON(200, gin.H{"message": "updated"})
}

func (c *NewsController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	c.service.Delete(uint(id))
	ctx.JSON(200, gin.H{"message": "deleted"})
}
