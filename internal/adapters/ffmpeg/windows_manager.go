package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/silvioubaldino/best-record-api/internal/core/domain"
	"github.com/silvioubaldino/best-record-api/internal/core/ports"
	"os"
	"os/exec"
	"strconv"
	"time"
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

	newffmpegStream.cmd = exec.Command("ffmpeg", "-f", "dshow", "-framerate", newffmpegStream.Fps, "-i", "video="+stream.CameraName,
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

	clipName := fmt.Sprintf("clip_%s_%s.mp4", newffmpegStream.cameraName, time.Now().Format("20060102_150405"))

	return extractClip(clipName, data)
}

func (w *windowsManager) GetAvailableCameras() (map[string]string, error) {
	//TODO ver como pego a lista de cameras disponiveis no windows
	cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")

	var out bytes.Buffer
	cmd.Stderr = &out

	_ = cmd.Run()

	output := out.String()

	//TODO ver como ficaria essa saída no windows e fazer a função de parse de acordo
	videoDevices := parseWindowsOutPut(output, "Windows video devices:")

	return videoDevices, nil
}

func parseWindowsOutPut(output, section string) map[string]string {

	panic(output)
}
