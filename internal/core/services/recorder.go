package services

import (
	"strings"

	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/ports"
)

type RecorderService struct {
	streamManager ports.StreamManager
	rgRepository  ports.RecordingGroupsRepository
}

func NewRecorderService(manager ports.StreamManager, repo ports.RecordingGroupsRepository) *RecorderService {
	return &RecorderService{
		streamManager: manager,
		rgRepository:  repo,
	}
}

func (s *RecorderService) GetRecordingGroups() ([]domain.RecordingGroup, error) {
	recGroup, err := s.rgRepository.GetRecordGroups()
	if err != nil {
		return nil, err
	}
	for i, group := range recGroup {
		for j, stream := range group.Streams {
			status, _ := s.streamManager.IsRecording(stream.ID)
			recGroup[i].Streams[j].IsRecording = status
		}
	}
	return recGroup, nil
}

func (s *RecorderService) StartGroupRecording(id uuid.UUID) error {
	recGroup, err := s.rgRepository.GetRecordGroup(id)
	if err != nil {
		return err
	}

	for _, stream := range recGroup.Streams {
		if err := s.streamManager.StartRecording(stream); err != nil {
			return err // TODO tratar erros para que tente iniciar todas as cameras mesmo que uma de erro
		}
	}
	return nil
}

func (s *RecorderService) StopRecording(id uuid.UUID) error {
	group, err := s.rgRepository.GetRecordGroup(id)
	if err != nil {
		return err
	}

	for _, stream := range group.Streams {
		if err := s.streamManager.StopRecording(stream.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *RecorderService) ClipRecording(id uuid.UUID, duration int) (string, error) {
	group, err := s.rgRepository.GetRecordGroup(id)
	if err != nil {
		return "", err
	}

	var clipNames []string
	for _, stream := range group.Streams {
		clip, err := s.streamManager.ClipRecording(stream.ID, duration)
		if err != nil {
			return "", err
		}
		clipNames = append(clipNames, clip)
	}
	return strings.Join(clipNames, ";"), nil
}

func (s *RecorderService) GetAvaiableCam() (map[string]string, error) {
	return s.streamManager.GetAvailableCameras()
}
