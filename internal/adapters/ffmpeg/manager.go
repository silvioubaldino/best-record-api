package ffmpeg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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

type ffmpegStream struct {
	id         uuid.UUID
	cameraName string
	VideoConfig
	circularBuffer *CircularBuffer
	mutex          sync.Mutex
	cmd            *exec.Cmd
}

type FFmpegManager struct {
	streams []*ffmpegStream
}

func NewFFmpegManager() ports.StreamManager {
	return &FFmpegManager{}
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

func (m *FFmpegManager) getStream(id uuid.UUID) (*ffmpegStream, error) {
	for _, stream := range m.streams {
		if id == stream.id {
			return stream, nil
		}
	}
	return nil, fmt.Errorf("incorrect stream ID")
}

func (m *FFmpegManager) addStream(stream domain.Stream) (*ffmpegStream, error) {
	existentStream, _ := m.getStream(stream.ID)
	if existentStream != nil {
		return existentStream, nil
	}

	newffmpegStream := toffmpegStream(stream)

	m.streams = append(m.streams, newffmpegStream)
	return newffmpegStream, nil
}

func (m *FFmpegManager) StartRecording(stream domain.Stream) error {
	newffmpegStream, err := m.addStream(stream)
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
	stream.Status = "recording"
	fmt.Printf("%s started", stream.CameraName)

	return nil
}

func (m *FFmpegManager) StopRecording(streamID uuid.UUID) error {
	newffmpegStream, err := m.getStream(streamID)
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

func (m *FFmpegManager) ClipRecording(streamID uuid.UUID, seconds int) (string, error) {
	newffmpegStream, err := m.getStream(streamID)
	if err != nil {
		return "", err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	// Read the last 'seconds' minutes from the circular buffer
	data, err := newffmpegStream.circularBuffer.ReadLastSeconds(seconds)
	if err != nil {
		return "", err
	}

	clipName := fmt.Sprintf("clip_%s_%s.mp4", newffmpegStream.cameraName, time.Now().Format(time.DateTime))
	return extractClip(clipName, data)

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

func (m *FFmpegManager) GetAvailableCameras() (map[string]string, error) {
	cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")

	var out bytes.Buffer
	cmd.Stderr = &out

	_ = cmd.Run()

	output := out.String()

	videoDevices := parseFFmpegOutput(output, "AVFoundation video devices:")

	return videoDevices, nil
}

func parseFFmpegOutput(output, section string) map[string]string {
	lines := strings.Split(output, "\n")
	var devices = make(map[string]string)
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
