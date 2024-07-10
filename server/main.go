package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/best-record-api/internal/adapters/controllers"
	"github.com/silvioubaldino/best-record-api/internal/adapters/ffmpeg"
	"github.com/silvioubaldino/best-record-api/internal/adapters/repositories"
	"github.com/silvioubaldino/best-record-api/internal/app"
	"github.com/silvioubaldino/best-record-api/internal/core/services"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permitir todas as origens
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	manager, err := ffmpeg.GetVideoManager()
	if err != nil {
		panic(err)
	}
	repo := repositories.NewTempoRepo()
	recorderService := services.NewRecorderService(manager, repo)

	recorderController := controllers.NewRecorderController(recorderService)

	app.SetupRoutes(r, recorderController)

	localIP := getLocalIP()
	if localIP == "" {
		fmt.Println("Não foi possível encontrar o endereço IP local.")
		os.Exit(1)
	}

	err = r.Run(localIP + ":8080")
	if err != nil {
		panic(err)
	}
}

func getLocalIP() string {
	environment := os.Getenv("environment")
	if environment == "development" {
		return "localhost"
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return ""
}
