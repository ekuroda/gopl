package main

import "io"

type counter struct {
	count  int64
	writer io.Writer
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := &counter{0, w}
	return io.Writer(c), &c.count
}

func (c *counter) Write(p []byte) (int, error) {
	count, err := c.writer.Write(p)
	c.count += int64(count)
	return count, err
}
