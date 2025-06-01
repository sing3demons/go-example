package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"github.com/sing3demons/go-example/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductController struct {
	db store.Storer
}

func NewProductController(db store.Storer) *ProductController {
	return &ProductController{db}
}

func (p *ProductController) Find(c *gin.Context) {
	products := []models.Product{}

	if err := p.db.Find(&products, bson.D{}); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			c.JSON(http.StatusNotFound, gin.H{
				"data": products,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})

}

func (p *ProductController) FindOne(c *gin.Context) {
	var product models.Product

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
		})
		return
	}

	if err := p.db.First(&product, bson.M{"_id": id}); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
			})
			return
		default:
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"data": product,
	})
}

func (p *ProductController) Create(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := p.db.Create(&product); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"data": product,
	})
}
