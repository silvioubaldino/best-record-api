package domain

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	_defaultOutputPath        = "Videos"
	_defaultRecordsFolderPath = "Best_Records"
)

type (
	Stream struct {
		ID          uuid.UUID
		CameraID    string
		CameraName  string
		IsRecording bool
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

	outputPath := filepath.Join(currentUser.HomeDir, _defaultOutputPath, _defaultRecordsFolderPath)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(outputPath, 0o755) // Permissões rwxr-xr-x
		if err != nil {
			return "", err
		}
	}
	return outputPath, nil
}
