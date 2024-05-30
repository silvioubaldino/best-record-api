package app

import (
	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/best-record-api/internal/adapters/controllers"
)

func SetupRoutes(r *gin.Engine) {
	recorderController := controllers.NewRecorderController()

	r.POST("/record", recorderController.StartRecording)
	r.POST("/stop", recorderController.StopRecording)
	r.POST("/clip", recorderController.ClipRecording)
	r.GET("/get-cameras", recorderController.GetAvailableCameras)
}
