package bzip

import (
	"io"
	"os/exec"
	"sync"
)

type writer struct {
	cmd *exec.Cmd
	w   io.WriteCloser
	sync.Mutex
}

// NewWriter ...
func NewWriter(out io.Writer) io.WriteCloser {
	cmd := exec.Command("/usr/bin/bzip2")
	cmd.Stdout = out
	w := &writer{cmd: cmd}
	return w
}

func (w *writer) Write(data []byte) (int, error) {
	w.Lock()
	defer func() {
		w.Unlock()
	}()

	if w.w == nil {
		stdin, err := w.cmd.StdinPipe()
		if err != nil {
			return 0, err
		}
		w.w = stdin
		if err := w.cmd.Start(); err != nil {
			return 0, err
		}
	}

	return w.w.Write(data)
}

func (w *writer) Close() error {
	w.Lock()
	defer func() {
		w.Unlock()
	}()

	if w.w == nil {
		return nil
	}
	if err := w.w.Close(); err != nil {
		return err
	}
	if err := w.cmd.Wait(); err != nil {
		return err
	}
	return nil
}
