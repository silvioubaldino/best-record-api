package domain

import (
	"github.com/google/uuid"
	"os"
	"os/user"
	"path/filepath"
)

type (
	Stream struct {
		ID          uuid.UUID
		CameraID    string
		CameraName  string
		Status      string
		Fps         string
		BitRate     int
		MaxDuration int
	}

	RecordingGroup struct {
		ID      uuid.UUID
		Name    string
		Streams []Stream
	}
)

func GetOutputPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(currentUser.HomeDir, "Videos", "Best Records")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0755) // Permissões rwxr-xr-x
		if err != nil {
			return "", err
		}
	}
	return outputPath, nil
}
