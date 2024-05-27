package main

import (
	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/best-record-api/internal/app"
)

func main() {
	r := gin.Default()
	app.SetupRoutes(r)
	r.Run(":8080")
}
