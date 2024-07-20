package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/ports"
)

type VideoConfig struct {
	Fps         string
	BitRate     int
	MaxDuration int
}

type Streams []*ffmpegStream

type ffmpegStream struct {
	id             uuid.UUID
	cameraName     string
	isRecording    bool
	circularBuffer *CircularBuffer
	mutex          sync.Mutex
	cmd            *exec.Cmd
	VideoConfig
}

func GetVideoManager() (ports.StreamManager, error) {
	so := runtime.GOOS
	switch so {
	case "darwin":
		return NewMacOSManager(), nil
	case "windows":
		return NewWindowsManager(), nil

	}

	return nil, fmt.Errorf("unsupported OS: %s", so)
}

func toffmpegStream(stream domain.Stream) *ffmpegStream {
	buffer := NewCircularBuffer(stream.BitRate, stream.MaxDuration)

	newffmpegStream := &ffmpegStream{
		id:          stream.ID,
		cameraName:  stream.CameraName,
		isRecording: stream.IsRecording,
		VideoConfig: VideoConfig{
			Fps:         stream.Fps,
			BitRate:     stream.BitRate,
			MaxDuration: stream.MaxDuration,
		},
		circularBuffer: buffer,
	}

	return newffmpegStream
}

func (m *Streams) getStream(id uuid.UUID) (*ffmpegStream, error) {
	for _, stream := range *m {
		if id == stream.id {
			return stream, nil
		}
	}
	return nil, fmt.Errorf("incorrect stream ID")
}

func (m *Streams) addStream(stream domain.Stream) (*ffmpegStream, error) {
	existentStream, _ := m.getStream(stream.ID)
	if existentStream != nil {
		return &ffmpegStream{}, fmt.Errorf("stream already exists")
	}

	newffmpegStream := toffmpegStream(stream)

	*m = append(*m, newffmpegStream)
	return newffmpegStream, nil
}

func (m *Streams) removeStream(id uuid.UUID) {
	for i, stream := range *m {
		if id == stream.id {
			*m = append((*m)[:i], (*m)[i+1:]...)
		}
	}
}

func extractClip(clipName string, data []byte) (string, error) {
	outputPath, err := domain.GetOutputPath()
	if err != nil {
		return "", err
	}

	outputPathFile := filepath.Join(outputPath, clipName)

	tempFile := fmt.Sprintf("temp_%d.ts", time.Now().Unix())
	fmt.Printf("%s", outputPathFile)
	file, err := os.Create(tempFile)
	if err != nil {
		return "", err
	}

	if _, err := file.Write(data); err != nil {
		return "", err
	}
	file.Close()

	ffmpegCmd := exec.Command("ffmpeg", "-i", tempFile, "-c:v", "libx264", "-preset", "fast", "-crf", "22", "-c:a", "aac", "-strict", "experimental", outputPathFile)
	ffmpegCmd.Stderr = os.Stderr
	if err := ffmpegCmd.Run(); err != nil {
		return "", err
	}

	if err := os.Remove(tempFile); err != nil {
		return "", err
	}
	return outputPathFile, nil
}
