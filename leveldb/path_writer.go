package leveldb

import (
	"encoding/binary"
	"log"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/syndtr/goleveldb/leveldb"
)

// PathWriter - buffered stdout writer with sync channel
type PathWriter struct {
	Conn      *Connection
	WaitGroup *sync.WaitGroup
	Queue     chan kv
}

// NewPathWriter - constructor
func NewPathWriter(conn *Connection) *PathWriter {
	w := &PathWriter{
		Conn:      conn,
		WaitGroup: &sync.WaitGroup{},
		Queue:     make(chan kv, batchSize*10),
	}

	// start writer routine
	w.WaitGroup.Add(1)
	go func() {
		batch := new(leveldb.Batch)
		for row := range w.Queue {

			// put
			batch.Put(row.Key, row.Val)

			// flush when full
			if batch.Len() >= batchSize {

				// write batch
				err := w.Conn.DB.Write(batch, nil)
				if err != nil {
					log.Println(err)
				}

				// reset batch
				batch.Reset()
			}
		}

		// write final batch
		err := w.Conn.DB.Write(batch, nil)
		if err != nil {
			log.Println(err)
		}

		w.WaitGroup.Done()
	}()

	return w
}

// Enqueue - close the channel and block until done
func (w *PathWriter) Enqueue(item *gosmparse.Way) {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(item.ID))

	// encoded path
	var value []byte

	// iterate over node refs, appending each int64 id to the value
	for _, ref := range item.NodeIDs {

		// encode id
		// @todo: use varint encoding to save bytes
		idBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idBytes, uint64(ref))

		// append to slice
		value = append(value, idBytes...)
	}

	w.Queue <- kv{Key: key, Val: value}
}

// Close - close the channel and block until done
func (w *PathWriter) Close() {
	close(w.Queue)
	w.WaitGroup.Wait()
}
