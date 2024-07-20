package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/services"
)

type RecorderController struct {
	service *services.RecorderService
}

func NewRecorderController(service *services.RecorderService) *RecorderController {
	return &RecorderController{
		service: service,
	}
}

func (r *RecorderController) GetRecordingGroups(c *gin.Context) {
	groups, err := r.service.GetRecordingGroups()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrGettingRecordingGroup.Error())
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (r *RecorderController) StartRecording(c *gin.Context) {
	idString := c.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	if err = r.service.StartGroupRecording(id); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrStartingRecording.Error())
		return
	}

	c.JSON(http.StatusOK, "Recording started")
}

func (r *RecorderController) StopRecording(c *gin.Context) {
	idString := c.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	if err = r.service.StopRecording(id); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrStopRecording.Error())
		return
	}

	c.JSON(http.StatusNoContent, "Recording stopped")
}

func (r *RecorderController) ClipRecording(c *gin.Context) {
	idString := c.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	var req struct {
		Duration int `json:"duration"`
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, ErrInvalidDuration.Error())
		return
	}

	path, err := r.service.ClipRecording(id, req.Duration)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrclipRecording.Error())
		return
	}

	paths := strings.Split(path, ";")
	filePath := strings.TrimSpace(paths[0])
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrFileNotFound.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"clipped_video_path": path})
}

func (r *RecorderController) GetAvailableCameras(c *gin.Context) {
	camList, err := r.service.GetAvaiableCam()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrGettingCameras.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"cameras": camList})
}

func (r *RecorderController) ServeClip(c *gin.Context) {
	fileName := c.Param("filename")

	outputPath, err := domain.GetOutputPath()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrGettingHomeDir.Error())
		return
	}

	path := filepath.Join(outputPath, fileName)
	remoteFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, ErrFileNotFound.Error())
		return
	}
	defer remoteFile.Close()

	c.FileAttachment(path, fileName)
}
