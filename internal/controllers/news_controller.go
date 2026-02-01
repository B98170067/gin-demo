package controller

import (
	model "gin-demo/internal/models"
	service "gin-demo/internal/services"
	errno "gin-demo/pkg/error"
	"gin-demo/pkg/response"
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
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	var status *int
	if s := ctx.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	data, total := c.service.GetPaged(page, size, status)
	ctx.JSON(200, gin.H{
		"data":  data,
		"total": total,
	})
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

// GetByID godoc
// @Summary Get news by id
// @Description Get single news
// @Tags News
// @Param id path int true "News ID"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /api/news/{id} [get]
func (c *NewsController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(errno.New(errno.ErrInvalidParam, "invalid id"))
		return
	}

	data, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.Error(errno.New(errno.ErrNotFound, "news not found"))
		return
	}

	response.Success(ctx, data)
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
