package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/controllers"
	"gorm.io/gorm"
)

func Router(r *gin.Engine, db *gorm.DB) {
	todoController := controllers.NewTodoController(db)

	r.GET("/todos", todoController.Index)
	r.POST("/todos", todoController.Create)

}
