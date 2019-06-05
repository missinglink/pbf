package leveldb

import (
	"encoding/binary"
	"log"
	"math"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/syndtr/goleveldb/leveldb"
)

var batchSize = 20000

// CoordWriter - buffered stdout writer with sync channel
type CoordWriter struct {
	Conn      *Connection
	WaitGroup *sync.WaitGroup
	Queue     chan kv
}

type kv struct {
	Key []byte
	Val []byte
}

// NewCoordWriter - constructor
func NewCoordWriter(conn *Connection) *CoordWriter {
	w := &CoordWriter{
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
func (w *CoordWriter) Enqueue(item *gosmparse.Node) {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(item.ID))

	// encode lat
	lat := make([]byte, 8)
	binary.BigEndian.PutUint64(lat, math.Float64bits(item.Lat))

	// encode lon
	lon := make([]byte, 8)
	binary.BigEndian.PutUint64(lon, math.Float64bits(item.Lon))

	// value
	value := append(lat, lon...)

	w.Queue <- kv{Key: key, Val: value}
}

// Close - close the channel and block until done
func (w *CoordWriter) Close() {
	close(w.Queue)
	w.WaitGroup.Wait()
}
