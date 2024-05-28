package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type VideoConfig struct {
	Fps         string
	BitRate     int
	MaxDuration int
}

type FFmpegManager struct {
	circularBuffer *CircularBuffer
	mutex          sync.Mutex
	ffmpegCmd      *exec.Cmd
	videoConfig    VideoConfig
}

func NewFFmpegManager(config VideoConfig) *FFmpegManager {
	return &FFmpegManager{
		circularBuffer: NewCircularBuffer(config.BitRate, config.MaxDuration),
		videoConfig:    config,
	}
}

func (m *FFmpegManager) StartRecording() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.ffmpegCmd = exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", m.videoConfig.Fps, "-i", "0",
		"-b:v", strconv.Itoa(m.videoConfig.BitRate)+"k", "-f", "mpegts", "pipe:1")
	m.ffmpegCmd.Stdout = m.circularBuffer
	m.ffmpegCmd.Stderr = os.Stderr
	if err := m.ffmpegCmd.Start(); err != nil {
		return err
	}

	return nil
}

func (m *FFmpegManager) StopRecording() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.ffmpegCmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}

	return nil
}

func (m *FFmpegManager) ClipRecording(seconds int) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Read the last 'seconds' minutes from the circular buffer
	data, err := m.circularBuffer.ReadLastSeconds(seconds)
	if err != nil {
		return "", err
	}

	return extractClip(data)

}

func extractClip(data []byte) (string, error) {
	// Write the data to a temporary TS file
	tempFile := fmt.Sprintf("temp_%d.ts", time.Now().Unix())
	file, err := os.Create(tempFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return "", err
	}

	// Convert the TS file to MP4
	output := fmt.Sprintf("clip_%d.mp4", time.Now().Unix())
	ffmpegCmd := exec.Command("ffmpeg", "-i", tempFile, "-c", "copy", output)
	ffmpegCmd.Stderr = os.Stderr
	if err := ffmpegCmd.Run(); err != nil {
		return "", err
	}

	// Delete the temporary TS file
	if err := os.Remove(tempFile); err != nil {
		return "", err
	}
	return output, nil
}
