package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProductRouter(r *gin.Engine, db *mongo.Collection) {
	productController := controllers.NewProductController(db)

	r.GET("/products", productController.Find)
	r.GET("/products/:id", productController.FindOne)
	r.POST("/products", productController.Create)
}
