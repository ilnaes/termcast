package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Recorder struct {
	writer []io.Writer
	start  time.Time
	dir    string
	pause  time.Time
	paused bool
	gap    time.Duration
}

func NewRecorder(f []io.Writer, dir string) *Recorder {
	return &Recorder{
		writer: f,
		dir:    dir,
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

	for _, w := range r.writer {
		_, err = io.WriteString(w, fmt.Sprintf("[%.2f, \"%s\", %s]\n", deltaTime, r.dir, b))
		if err != nil {
			return 0, err
		}
	}

	return len(p), err
}
