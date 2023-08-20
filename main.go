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

	db, err := db.ConnectDB()

	if err != nil {
		panic(err)
	}

	r := gin.Default()

	router.Router(r, db)

	r.Run(":" + os.Getenv("PORT"))

}
