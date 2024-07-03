package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/sftp"

	"github.com/silvioubaldino/best-record-api/internal/core/services"
)

type RecorderController struct {
	service    *services.RecorderService
	sftpClient *sftp.Client
}

func NewRecorderController(service *services.RecorderService, sftp *sftp.Client) *RecorderController {
	return &RecorderController{
		service:    service,
		sftpClient: sftp,
	}
}

func (r *RecorderController) GetRecordingGroups(c *gin.Context) {
	groups, err := r.service.GetRecordingGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recording_groups": groups})
}

func (r *RecorderController) StartRecording(c *gin.Context) {
	idString := c.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.service.StartGroupRecording(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recording started"})
}

func (r *RecorderController) StopRecording(c *gin.Context) {
	idString := c.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.service.StopRecording(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recording stopped"})
}

func (r *RecorderController) ClipRecording(c *gin.Context) {
	idString := c.Query("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Duration int `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path, err := r.service.ClipRecording(id, req.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paths := strings.Split(path, ";")
	filePath := strings.TrimSpace(paths[0])
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clipped_video_path": path})
}

func (r *RecorderController) GetAvailableCameras(c *gin.Context) {
	camList, err := r.service.GetAvaiableCam()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cameras": camList})
}

func (r *RecorderController) ServeClip(c *gin.Context) {
	fileName := c.Param("filename")

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user: ", err)
	}
	path := filepath.Join(currentUser.HomeDir, "Videos", fileName)
	remoteFile, err := r.sftpClient.Open(path + fileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	defer remoteFile.Close()

	_, err = io.Copy(c.Writer, remoteFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clip downloaded"})
}
