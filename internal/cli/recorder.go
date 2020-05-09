package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

type Recorder struct {
	writer io.Writer

	start  time.Time
	pause  time.Time
	paused bool
	gap    time.Duration

	direction string
	*sync.Mutex
}

func NewRecorder(writer io.Writer, dir string) *Recorder {
	return &Recorder{
		writer:    writer,
		start:     time.Now(),
		direction: dir,
		Mutex:     &sync.Mutex{},
	}
}

func (r *Recorder) Pause() {
	// TODO
}

func (r *Recorder) Resume() {
	// TODO
}

var header string = `{"version": 2, "width": %d, "height": %d, "timestamp": %d, "env": {"SHELL": "%s", "TERM": "%s"}}
`

func (r *Recorder) WriteHead(cmd string) (err error) {
	w, h, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	_, err = r.writer.Write([]byte(fmt.Sprintf(header, w, h, time.Now().Unix(), cmd, os.Getenv("TERM"))))
	return err
}

func (r *Recorder) Write(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()

	if r.paused {
		return 0, nil
	}

	deltaTime := float32(time.Since(r.start).Milliseconds()-r.gap.Milliseconds()) / 1000.

	b, err := json.Marshal(string(p))
	if err != nil {
		return 0, err
	}

	_, err = io.WriteString(r.writer, fmt.Sprintf("[%.2f, \"%s\", %s]\n",
		deltaTime, r.direction, b))

	if err != nil {
		return 0, err
	}

	return len(p), err
}
