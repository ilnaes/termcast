package cli

import (
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

func rec(wr io.Writer) error {
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

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/rec", nil)
	if err != nil {
		return err
	}
	ws := WsWriter{Conn: conn}

	outMwr := io.MultiWriter(os.Stdout, NewRecorder([]io.Writer{wr, ws}, "o"))

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
