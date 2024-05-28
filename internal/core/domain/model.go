package domain

type (
	Stream struct {
		ID          string
		cameraName  string
		Status      string
		Fps         string
		BitRate     string
		MaxDuration int
	}

	RecordingGroup struct {
		ID      string
		Name    string
		Streams []Stream
	}
)
