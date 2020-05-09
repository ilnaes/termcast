package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	. "github.com/ilnaes/termcast/internal"
	"golang.org/x/crypto/ssh/terminal"
)

func getCmd() string {
	return os.Getenv("SHELL")
}

func rec(w io.Writer) error {
	// Create arbitrary command.
	cmd := getCmd()
	c := exec.Command(cmd)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
	fmt.Printf("\x1b[0;0H\x1b[2J")

	wr := NewRecorder(w, "o")
	wr.WriteHead(cmd)
	writers := []io.Writer{os.Stdout, wr}

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/rec", nil)
	if err == nil {
		ws := NewRecorder(WsWriter{Conn: conn}, "o")
		writers = append(writers, ws)

		fmt.Printf("\x1b[32mConnected to server\x1b[00m\r\n")
	}

	outMwr := io.MultiWriter(writers...)

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(outMwr, ptmx)

	return nil
}

func Run() {
	f, _ := os.Create("output.cast")
	if err := rec(f); err != nil {
		log.Fatal(err)
	}
}
