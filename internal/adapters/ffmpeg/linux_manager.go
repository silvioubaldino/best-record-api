package ffmpeg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/ports"
)

type linuxManager struct {
	streams Streams
}

func NewLinuxManager() ports.StreamManager {
	return &linuxManager{}
}

func (m *linuxManager) StartRecording(stream domain.Stream) error {
	newffmpegStream, err := m.streams.addStream(stream)
	if err != nil {
		return err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	streamURL := "rtmp://nginx:1935/live/nome_do_stream"
	newffmpegStream.cmd = exec.Command("ffmpeg", "-i", streamURL, "-f", "mpegts", "pipe:1")
	newffmpegStream.cmd.Stdout = newffmpegStream.circularBuffer
	newffmpegStream.cmd.Stderr = os.Stderr
	if err := newffmpegStream.cmd.Start(); err != nil {
		return err
	}
	newffmpegStream.isRecording = true
	fmt.Printf("%s started", stream.CameraName)

	return nil
}

func (m *linuxManager) StopRecording(streamID uuid.UUID) error {
	newffmpegStream, err := m.streams.getStream(streamID)
	if err != nil {
		return err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	if err := newffmpegStream.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	m.streams.removeStream(streamID)
	return nil
}

func (m *linuxManager) ClipRecording(streamID uuid.UUID, seconds int) (string, error) {
	newffmpegStream, err := m.streams.getStream(streamID)
	if err != nil {
		return "", err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	data, err := newffmpegStream.circularBuffer.ReadLastSeconds(seconds)
	if err != nil {
		return "", err
	}

	clipName := fmt.Sprintf("clip_%s_%s.mp4", newffmpegStream.cameraName, time.Now().Format(time.DateTime))
	return extractClip(clipName, data)
}

func (m *linuxManager) GetAvailableCameras() (map[string]string, error) {
	cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")

	var out bytes.Buffer
	cmd.Stderr = &out

	_ = cmd.Run()

	output := out.String()

	videoDevices := parseAvFoundationOutput(output, "AVFoundation video devices:")

	return videoDevices, nil
}

func (m *linuxManager) IsRecording(streamID uuid.UUID) (bool, error) {
	newffmpegStream, err := m.streams.getStream(streamID)
	if err != nil {
		return false, err
	}
	return newffmpegStream.isRecording, nil
}
