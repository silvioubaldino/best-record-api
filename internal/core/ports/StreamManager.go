package ports

import (
	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
)

type StreamManager interface {
	StartRecording(stream domain.Stream) error
	StopRecording(streamID uuid.UUID) error
	ClipRecording(streamID uuid.UUID, seconds int) (string, error)
}
