package ports

import (
	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
)

type RecordingGroupsRepository interface {
	GetRecordGroups() ([]domain.RecordingGroup, error)
	GetRecordGroup(id uuid.UUID) (domain.RecordingGroup, error)
}
