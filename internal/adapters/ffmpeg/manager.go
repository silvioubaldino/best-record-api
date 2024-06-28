package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
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
	id         uuid.UUID
	cameraName string
	VideoConfig
	circularBuffer *CircularBuffer
	mutex          sync.Mutex
	cmd            *exec.Cmd
}

type FFmpegManager struct {
	streams Streams
}

func GetVideoManager() (ports.StreamManager, error) {
	so := runtime.GOOS
	switch so {
	case "darwin":
		return NewMacOSManager(), nil
	case "linux":
		return NewLinuxManager(), nil
	case "windows":
		return nil, nil

	}

	return nil, fmt.Errorf("unsupported OS: %s", so)
}

func toffmpegStream(stream domain.Stream) *ffmpegStream {
	buffer := NewCircularBuffer(stream.BitRate, stream.MaxDuration)

	newffmpegStream := &ffmpegStream{
		id:         stream.ID,
		cameraName: stream.CameraName,
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
		return existentStream, nil
	}

	newffmpegStream := toffmpegStream(stream)

	*m = append(*m, newffmpegStream)
	return newffmpegStream, nil
}

func extractClip(clipName string, data []byte) (string, error) {
	tempFile := fmt.Sprintf("temp_%d.ts", time.Now().Unix())
	file, err := os.Create(tempFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return "", err
	}

	ffmpegCmd := exec.Command("ffmpeg", "-i", tempFile, "-c", "copy", clipName)
	ffmpegCmd.Stderr = os.Stderr
	if err := ffmpegCmd.Run(); err != nil {
		return "", err
	}

	if err := os.Remove(tempFile); err != nil {
		return "", err
	}
	return clipName, nil
}
