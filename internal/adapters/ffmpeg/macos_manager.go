package ffmpeg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/ports"
)

type macOSManager struct {
	streams Streams
}

func NewMacOSManager() ports.StreamManager {
	return &macOSManager{}
}

func (m *macOSManager) StartRecording(stream domain.Stream) error {
	newffmpegStream, err := m.streams.addStream(stream)
	if err != nil {
		return err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	newffmpegStream.cmd = exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", newffmpegStream.Fps, "-i", stream.CameraID,
		"-b:v", strconv.Itoa(newffmpegStream.BitRate)+"k", "-f", "mpegts", "pipe:1")
	newffmpegStream.cmd.Stdout = newffmpegStream.circularBuffer
	newffmpegStream.cmd.Stderr = os.Stderr
	if err := newffmpegStream.cmd.Start(); err != nil {
		return err
	}
	newffmpegStream.isRecording = true
	fmt.Printf("%s started", stream.CameraName)

	return nil
}

func (m *macOSManager) StopRecording(streamID uuid.UUID) error {
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

func (m *macOSManager) ClipRecording(streamID uuid.UUID, seconds int) (string, error) {
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

func (m *macOSManager) GetAvailableCameras() (map[string]string, error) {
	cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")

	var out bytes.Buffer
	cmd.Stderr = &out

	_ = cmd.Run()

	output := out.String()

	videoDevices := parseAvFoundationOutput(output, "AVFoundation video devices:")

	return videoDevices, nil
}

func (m *macOSManager) IsRecording(streamID uuid.UUID) (bool, error) {
	newffmpegStream, err := m.streams.getStream(streamID)
	if err != nil {
		return false, err
	}
	return newffmpegStream.isRecording, nil
}

func parseAvFoundationOutput(output, section string) map[string]string {
	lines := strings.Split(output, "\n")
	devices := make(map[string]string)
	var capture bool

	for _, line := range lines {
		if strings.Contains(line, section) {
			capture = true
			continue
		}
		if capture {
			if strings.Contains(line, "AVFoundation") && strings.Contains(line, "devices:") {
				break
			}
			re := regexp.MustCompile(`\[(\d+)\] (.+)`)
			match := re.FindStringSubmatch(line)
			if len(match) == 3 {
				devices[match[1]] = match[2]
			}
		}
	}
	return devices
}
