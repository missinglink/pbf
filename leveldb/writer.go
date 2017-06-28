package leveldb

import (
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var batchSize = 10000

// CoordWriter - buffered stdout writer with sync channel
type CoordWriter struct {
	Conn      *Connection
	WaitGroup *sync.WaitGroup
	Queue     chan []byte
}

// NewCoordWriter - constructor
func NewCoordWriter(conn *Connection) *CoordWriter {
	w := &CoordWriter{
		Conn:      conn,
		WaitGroup: &sync.WaitGroup{},
		Queue:     make(chan []byte, batchSize*2),
	}

	// start writer routine
	w.WaitGroup.Add(1)
	go func() {
		batch := new(leveldb.Batch)
		for encoded := range w.Queue {
			batch.Put(encoded[:9], encoded[9:])
			if batch.Len() >= batchSize {
				err := w.Conn.DB.Write(batch, nil)
				if err != nil {
					log.Println(err)
				}
				batch.Reset()
			}
		}
		err := w.Conn.DB.Write(batch, nil)
		if err != nil {
			log.Println(err)
		}
		w.WaitGroup.Done()
	}()

	return w
}

// Close - close the channel and block until done
func (w *CoordWriter) Close() {
	close(w.Queue)
	w.WaitGroup.Wait()
	w.Conn.Close()
}
