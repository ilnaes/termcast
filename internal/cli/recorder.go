package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Recorder struct {
	writer io.Writer
	start  time.Time
	pause  time.Time
	paused bool
	gap    time.Duration
}

func NewRecorder(f io.Writer) *Recorder {
	return &Recorder{
		writer: f,
		start:  time.Now(),
	}
}

func (r *Recorder) Pause() {
	// TODO
}

func (r *Recorder) Resume() {
	// TODO
}

func (r *Recorder) WriteHead() (err error) {
	// TODO
	return nil
}

func (r *Recorder) Write(p []byte) (n int, err error) {
	if r.paused {
		return 0, nil
	}

	deltaTime := float32(time.Since(r.start).Milliseconds()-r.gap.Milliseconds()) / 1000.

	b, err := json.Marshal(string(p))
	if err != nil {
		return 0, err
	}

	_, err = io.WriteString(r.writer, fmt.Sprintf("[%.2f, %s]\n", deltaTime, b))
	if err != nil {
		return 0, err
	}

	return len(p), err
}
