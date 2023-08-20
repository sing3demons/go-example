package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"gorm.io/gorm"
)

type TodoController struct {
	db *gorm.DB
}

func NewTodoController(db *gorm.DB) *TodoController {
	return &TodoController{db}
}

func (t *TodoController) Index(c *gin.Context) {
	todos := []models.Todo{}

	if err := t.db.Find(&todos).Error; err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}

type TodoCreateRequest struct {
	Title string `json:"title" binding:"required"`
}

func (t *TodoController) Create(c *gin.Context) {
	var req TodoCreateRequest
	c.ShouldBindJSON(&req)

	todo := models.Todo{
		Title:     req.Title,
		Completed: false,
	}

	t.db.Create(&todo)

	c.JSON(http.StatusCreated, gin.H{
		"data": todo,
	})
}
