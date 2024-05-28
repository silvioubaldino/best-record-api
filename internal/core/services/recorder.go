package services

import (
	"github.com/silvioubaldino/best-record-api/internal/adapters/ffmpeg"
)

type RecorderService struct {
	manager *ffmpeg.FFmpegManager
}

func NewRecorderService() *RecorderService {
	return &RecorderService{
		manager: ffmpeg.NewFFmpegManager(ffmpeg.VideoConfig{
			Fps:         "30",
			BitRate:     8000,
			MaxDuration: 30 * 60,
		}),
	}
}

func (s *RecorderService) StartRecording() error {

	return s.manager.StartRecording()
}

func (s *RecorderService) StopRecording() error {
	return s.manager.StopRecording()
}

func (s *RecorderService) ClipRecording(duration int) (string, error) {
	return s.manager.ClipRecording(duration)
}
