package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/ports"
)

type tempoRepo struct{}

func NewTempoRepo() ports.RecordingGroupsRepository {
	return tempoRepo{}
}

func (t tempoRepo) GetRecordGroups() ([]domain.RecordingGroup, error) {
	recordingGroups, err := readJSON("./config.json")
	if err != nil {
		return []domain.RecordingGroup{}, err
	}
	return recordingGroups, nil
}

func (t tempoRepo) GetRecordGroup(id uuid.UUID) (domain.RecordingGroup, error) {
	groups, err := t.GetRecordGroups()
	if err != nil {
		return domain.RecordingGroup{}, nil
	}
	for _, group := range groups {
		if group.ID == id {
			return group, nil
		}
	}
	return domain.RecordingGroup{}, fmt.Errorf("incorrect record group ID")
}

func readJSON(config string) ([]domain.RecordingGroup, error) {
	var recordingGroups []domain.RecordingGroup

	file, err := os.Open(config)
	if err != nil {
		return recordingGroups, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return recordingGroups, err
	}

	err = json.Unmarshal(bytes, &recordingGroups)
	if err != nil {
		return recordingGroups, err
	}

	return recordingGroups, nil
}
