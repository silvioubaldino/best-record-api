package main

import (
	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/best-record-api/internal/app"
)

func main() {
	r := gin.Default()
	err := app.SetupRoutes(r)
	if err != nil {
		panic(err)
	}
	r.Run(":8080")
}
