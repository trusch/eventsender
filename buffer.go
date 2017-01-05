package eventsender

import (
	"bytes"
	"io"
)

// Buffer is like a normal bytes.Buffer but blocks on Reads instead of returning EOF
// It will only return EOF if Close is called.
type Buffer struct {
	ch       chan []byte
	isClosed bool
	buffered *bytes.Buffer
}

// NewBuffer creates a new Buffer
func NewBuffer() *Buffer {
	return &Buffer{make(chan []byte, 32), false, &bytes.Buffer{}}
}

// Write feeds p into the buffer
func (buffer *Buffer) Write(p []byte) (n int, err error) {
	if buffer.isClosed {
		return 0, io.EOF
	}
	buffer.ch <- p
	return len(p), nil
}

// Read reads up to len(p) bytes into p.
// It will block until data is available or Close is called
func (buffer *Buffer) Read(p []byte) (n int, err error) {
	if buffer.buffered.Len() > 0 {
		n, err := buffer.buffered.Read(p)
		if err == io.EOF {
			return n, nil
		}
		return n, err
	}
	if data, ok := <-buffer.ch; ok {
		buffer.buffered.Write(data)
		return buffer.Read(p)
	}
	return 0, io.EOF
}

// Close closes the buffer, so that EOF will be returned on Reads
func (buffer *Buffer) Close() error {
	if !buffer.isClosed {
		close(buffer.ch)
		buffer.isClosed = true
	}
	return nil
}
