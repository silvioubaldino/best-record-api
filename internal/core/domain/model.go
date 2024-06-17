package domain

import (
	"github.com/google/uuid"
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
