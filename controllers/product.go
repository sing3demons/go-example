package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductController struct {
	db *mongo.Collection
}

func NewProductController(db *mongo.Collection) *ProductController {
	return &ProductController{db}
}

func (p *ProductController) Find(c *gin.Context) {
	products := []models.Product{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := p.db.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &products); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"data": products,
	})

}

func (p *ProductController) FindOne(c *gin.Context) {
	var product models.Product

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	if err := p.db.FindOne(ctx, bson.M{"_id": id}).Decode(&product); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"data": product,
	})
}

func (p *ProductController) Create(c *gin.Context) {
	var product models.Product
	c.ShouldBindJSON(&product)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := p.db.InsertOne(ctx, product)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	product.ID = r.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{
		"data": product,
	})
}
