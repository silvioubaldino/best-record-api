package controllers

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/best-record-api/internal/core/services"
)

type RecorderController struct {
	service *services.RecorderService
}

func NewRecorderController() (*RecorderController, error) {
	service, err := services.NewRecorderService()
	if err != nil {
		return nil, err
	}
	return &RecorderController{service: service}, nil
}

func (c *RecorderController) StartRecording(ctx *gin.Context) {
	idString := ctx.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.StartGroupRecording(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recording started"})
}

func (c *RecorderController) StopRecording(ctx *gin.Context) {
	idString := ctx.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.StopRecording(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recording stopped"})
}

func (c *RecorderController) ClipRecording(ctx *gin.Context) {
	idString := ctx.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Duration int `json:"duration"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path, err := c.service.ClipRecording(id, req.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"clipped_video_path": path})
}

func (c RecorderController) GetAvailableCameras(ctx *gin.Context) {
	camList, err := c.service.GetAvaiableCam()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"cameras": camList})
}
