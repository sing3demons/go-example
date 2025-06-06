package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/controllers"
	"github.com/sing3demons/go-example/store"
)

func Router(r *gin.Engine, db store.Storer) {
	todoController := controllers.NewTodoController(db)

	r.GET("/todos", todoController.Index)
	r.POST("/todos", todoController.Create)

}

func ProductRouter(r *gin.Engine, db store.Storer) {
	productController := controllers.NewProductController(db)

	r.GET("/products", productController.Find)
	r.GET("/products/:id", productController.FindOne)
	r.POST("/products", productController.Create)
}
