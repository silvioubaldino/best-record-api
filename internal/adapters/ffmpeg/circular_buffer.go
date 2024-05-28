package ffmpeg

import "errors"

type CircularBuffer struct {
	data []byte
	max  int
}

func NewCircularBuffer(bitRate int, maxSeconds int) *CircularBuffer {
	bytesPerSecond := bitRate * 1024 / 8

	maxBytes := maxSeconds * bytesPerSecond

	return &CircularBuffer{
		data: make([]byte, 0),
		max:  maxBytes,
	}
}

func (b *CircularBuffer) Write(p []byte) (n int, err error) {
	n = len(p)
	overflow := len(b.data) + n - b.max

	if overflow > 0 {
		b.data = b.data[overflow:]
	}

	b.data = append(b.data, p...)

	return n, nil
}

func (b *CircularBuffer) ReadLastSeconds(n int) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("n must be non-negative")
	}

	// Calculate the number of bytes per second for a bitrate of 5000 kbps
	bytesPerSecond := 8000 * 1024 / 8

	// Calculate the number of bytes to read
	bytesToRead := n * bytesPerSecond

	start := len(b.data) - bytesToRead
	if start < 0 {
		start = 0
	}

	return b.data[start:], nil
}
