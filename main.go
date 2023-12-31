package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/go-example/db"
	"github.com/sing3demons/go-example/router"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}

	gorm, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	mongo, err := db.ConnectMonoDB()

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	router.Router(r, gorm)
	router.ProductRouter(r, mongo)

	r.Run(":" + os.Getenv("PORT"))

}
