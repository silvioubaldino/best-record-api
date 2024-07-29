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

type windowsManager struct {
	streams Streams
}

func NewWindowsManager() ports.StreamManager {
	return &windowsManager{}
}

func (w *windowsManager) StartRecording(stream domain.Stream) error {
	newffmpegStream, err := w.streams.addStream(stream)
	if err != nil {
		return err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	fmt.Printf("%s", stream.ID)

	newffmpegStream.cmd = exec.Command("ffmpeg", "-i", stream.CameraName,
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

func (w *windowsManager) StopRecording(streamID uuid.UUID) error {
	newffmpegStream, err := w.streams.getStream(streamID)
	if err != nil {
		return err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	if err := newffmpegStream.cmd.Process.Kill(); err != nil {
		return err
	}
	w.streams.removeStream(streamID)

	return nil
}

func (w *windowsManager) ClipRecording(streamID uuid.UUID, seconds int) (string, error) {
	newffmpegStream, err := w.streams.getStream(streamID)
	if err != nil {
		return "", err
	}

	newffmpegStream.mutex.Lock()
	defer newffmpegStream.mutex.Unlock()

	data, err := newffmpegStream.circularBuffer.ReadLastSeconds(seconds)
	if err != nil {
		return "", err
	}

	clipName := fmt.Sprintf("clip_%s_%s.mp4", newffmpegStream.id, time.Now().Format("20060102_150405"))

	return extractClip(clipName, data)
}

func (w *windowsManager) GetAvailableCameras() (map[string]string, error) {
	cmd := exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")

	var out bytes.Buffer
	cmd.Stderr = &out

	_ = cmd.Run()

	output := out.String()

	videoDevices := parseWindowsOutPut(output)

	fmt.Printf("%s", videoDevices)

	return videoDevices, nil
}

func (w *windowsManager) IsRecording(streamID uuid.UUID) (bool, error) {
	newffmpegStream, err := w.streams.getStream(streamID)
	if err != nil {
		return false, err
	}
	return newffmpegStream.isRecording, nil
}

func parseWindowsOutPut(output string) map[string]string {
	lines := strings.Split(output, "\n")
	devices := make(map[string]string)

	var currentDeviceName string
	re := regexp.MustCompile(`"(.+?)" \(video\)`)
	reAlt := regexp.MustCompile(`Alternative name "(.+?)"`)

	for _, line := range lines {
		if match := re.FindStringSubmatch(line); len(match) == 2 {
			currentDeviceName = match[1]
		} else if match := reAlt.FindStringSubmatch(line); len(match) == 2 && currentDeviceName != "" {
			devices[match[1]] = currentDeviceName
			currentDeviceName = ""
		}
	}
	return devices
}
