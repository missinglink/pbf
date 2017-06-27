package lib

import (
	"bufio"
	"os"
	"sync"
)

// BufferedWriter - buffered stdout writer with sync channel
type BufferedWriter struct {
	writer    *bufio.Writer
	waitGroup *sync.WaitGroup
	Queue     chan []byte
}

// NewBufferedWriter - constructor
func NewBufferedWriter() *BufferedWriter {
	w := &BufferedWriter{
		writer:    bufio.NewWriter(os.Stdout),
		waitGroup: &sync.WaitGroup{},
		Queue:     make(chan []byte, 10000),
	}

	// start writer routine
	w.waitGroup.Add(1)
	go func() {
		for bytes := range w.Queue {
			w.writer.Write(bytes)
			w.writer.WriteRune('\n')
		}
		w.writer.Flush()
		w.waitGroup.Done()
	}()

	return w
}

// Close - close the channel and block until done
func (w *BufferedWriter) Close() {
	close(w.Queue)
	w.waitGroup.Wait()
}
