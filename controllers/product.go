package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/go-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductController struct {
	db *mongo.Collection
}

func NewProductController(db *mongo.Collection) *ProductController {
	return &ProductController{db}
}

func (p *ProductController) Find(c *gin.Context) {
	products := []models.Product{}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	search := c.Query("search")
	opts := options.Find().SetLimit(int64(limit))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if search != "" {
		filter = bson.M{"name": bson.M{"$regex": search}}
	}

	cursor, err := p.db.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		product := models.Product{}
		if err := cursor.Decode(&product); err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		products = append(products, product)
	}

	total, err := p.db.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"total": total,
		"data":  products,
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
		return
	}

	if err := p.db.FindOne(ctx, bson.M{"_id": id}).Decode(&product); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
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
		return
	}

	product.ID = r.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{
		"data": product,
	})
}
