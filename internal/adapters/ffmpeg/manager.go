package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
)

type FFmpegManager struct {
	currentRecording *domain.Recording
	mutex            sync.Mutex
	ffmpegCmd        *exec.Cmd
}

func NewFFmpegManager() *FFmpegManager {
	return &FFmpegManager{}
}

func (m *FFmpegManager) StartRecording(input, output string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentRecording != nil && m.currentRecording.Status == "recording" {
		return errors.New("already recording")
	}

	m.ffmpegCmd = exec.Command("ffmpeg", "-i", input, "-c", "copy", output)
	if err := m.ffmpegCmd.Start(); err != nil {
		return err
	}

	m.currentRecording = &domain.Recording{
		StartTime: time.Now(),
		Status:    "recording",
		FilePath:  output,
	}

	return nil
}

func (m *FFmpegManager) StopRecording() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentRecording == nil || m.currentRecording.Status != "recording" {
		return errors.New("not currently recording")
	}

	if err := m.ffmpegCmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}

	m.currentRecording.EndTime = time.Now()
	m.currentRecording.Status = "stopped"

	return nil
}

func (m *FFmpegManager) GetStatus() (domain.Recording, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentRecording == nil {
		return domain.Recording{}, errors.New("no current recording")
	}

	return *m.currentRecording, nil
}

func (m *FFmpegManager) ClipRecording(output string, duration int) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentRecording == nil || m.currentRecording.Status != "stopped" {
		return "", errors.New("no recording to clip")
	}

	cmd := exec.Command("ffmpeg", "-sseof", fmt.Sprintf("-%d", duration), "-i", m.currentRecording.FilePath, "-c", "copy", output)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %s", string(outputBytes))
	}

	return output, nil
}
