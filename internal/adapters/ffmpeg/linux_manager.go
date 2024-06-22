package ffmpeg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	newffmpegStream.cmd = exec.Command("ffmpeg", "-f", "v4l2", "-framerate", newffmpegStream.Fps, "-i", stream.CameraID,
		"-b:v", strconv.Itoa(newffmpegStream.BitRate)+"k", "-f", "mpegts", "pipe:1")
	newffmpegStream.cmd.Stdout = newffmpegStream.circularBuffer
	newffmpegStream.cmd.Stderr = os.Stderr
	if err := newffmpegStream.cmd.Start(); err != nil {
		return err
	}
	stream.Status = "recording"
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
	cmd := exec.Command("v4l2-ctl", "--list-devices")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar comando v4l2-ctl: %w", err)
	}

	output := out.String()
	return parseV4L2Output(output), nil
}

func parseV4L2Output(output string) map[string]string {
	lines := strings.Split(output, "\n")
	devices := make(map[string]string)
	var currentDevice string

	for _, line := range lines {
		if strings.HasSuffix(line, ":") {
			currentDevice = strings.TrimSuffix(line, ":")
		} else if strings.HasPrefix(line, "/dev/") {
			devices[line] = currentDevice
		}
	}

	return devices
}
