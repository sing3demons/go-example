package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"github.com/sing3demons/go-example/store"
	"gorm.io/gorm"
)

type TodoController struct {
	db store.Storer
}

func NewTodoController(db store.Storer) *TodoController {
	return &TodoController{db}
}

func (t *TodoController) Index(c *gin.Context) {
	todos := []models.Todo{}

	if err := t.db.Find(&todos); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"message": "No todos found",
			})
			return
		}

		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	todo := models.Todo{
		Title:     req.Title,
		Completed: false,
	}

	if err := t.db.Create(&todo); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": todo,
	})
}
