package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/best-record-api/internal/core/services"
)

type RecorderController struct {
	service *services.RecorderService
}

func NewRecorderController() *RecorderController {
	service := services.NewRecorderService()
	return &RecorderController{service: service}
}

func (c *RecorderController) StartRecording(ctx *gin.Context) {
	input := ctx.Query("input")
	output := ctx.Query("output")

	if err := c.service.StartRecording(input, output); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recording started"})
}

func (c *RecorderController) StopRecording(ctx *gin.Context) {
	if err := c.service.StopRecording(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recording stopped"})
}

func (c *RecorderController) GetStatus(ctx *gin.Context) {
	status, err := c.service.GetStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": status})
}

func (c *RecorderController) ClipRecording(ctx *gin.Context) {
	var req struct {
		Output   string `json:"output"`
		Duration int    `json:"duration"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path, err := c.service.ClipRecording(req.Output, req.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"clipped_video_path": path})
}
