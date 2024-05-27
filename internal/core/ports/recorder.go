package ports

import "github.com/silvioubaldino/best-record-api/internal/core/domain"

type Recorder interface {
	StartRecording(input, output string) error
	StopRecording() error
	GetStatus() (domain.Recording, error)
	ClipRecording(output string, duration int) (string, error)
}
