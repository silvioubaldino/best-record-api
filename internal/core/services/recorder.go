package services

import (
	"github.com/silvioubaldino/best-record-api/internal/adapters/ffmpeg"
	"github.com/silvioubaldino/best-record-api/internal/core/domain"
)

type RecorderService struct {
	manager *ffmpeg.FFmpegManager
}

func NewRecorderService() *RecorderService {
	return &RecorderService{
		manager: ffmpeg.NewFFmpegManager(),
	}
}

func (s *RecorderService) StartRecording(input, output string) error {
	return s.manager.StartRecording(input, output)
}

func (s *RecorderService) StopRecording() error {
	return s.manager.StopRecording()
}

func (s *RecorderService) GetStatus() (domain.Recording, error) {
	return s.manager.GetStatus()
}

func (s *RecorderService) ClipRecording(output string, duration int) (string, error) {
	return s.manager.ClipRecording(output, duration)
}
