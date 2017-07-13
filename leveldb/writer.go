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

// WriteQueue - a channel + waitgroup for when it's done
type WriteQueue struct {
	Conn      *Connection
	WaitGroup *sync.WaitGroup
	Chan      chan kv
}

// Start the queue
func (q *WriteQueue) Start() {

	// start writer routine
	q.WaitGroup.Add(1)
	go func() {
		batch := new(leveldb.Batch)
		for row := range q.Chan {

			// put
			batch.Put(row.Key, row.Val)

			// flush when full
			if batch.Len() >= batchSize {

				// write batch
				err := q.Conn.DB.Write(batch, nil)
				if err != nil {
					log.Println(err)
				}

				// reset batch
				batch.Reset()
			}
		}

		// write final batch
		err := q.Conn.DB.Write(batch, nil)
		if err != nil {
			log.Println(err)
		}

		q.WaitGroup.Done()
	}()
}

// Close - close the channel and block until done
func (q *WriteQueue) Close() {
	close(q.Chan)
	q.WaitGroup.Wait()
}

// Writer - buffered stdout writer with sync channel
type Writer struct {
	Conn      *Connection
	NodeQueue *WriteQueue
	WayQueue  *WriteQueue
}

type kv struct {
	Key []byte
	Val []byte
}

// NewWriter - constructor
func NewWriter(conn *Connection) *Writer {
	var w = &Writer{
		Conn: conn,
		NodeQueue: &WriteQueue{
			Conn:      conn,
			WaitGroup: &sync.WaitGroup{},
			Chan:      make(chan kv, batchSize*10),
		},
		WayQueue: &WriteQueue{
			Conn:      conn,
			WaitGroup: &sync.WaitGroup{},
			Chan:      make(chan kv, batchSize*10),
		},
	}

	w.NodeQueue.Start()
	w.WayQueue.Start()

	return w
}

// EnqueueNode - enqueue node bytes to be saved to db
func (w *Writer) EnqueueNode(item *gosmparse.Node) {

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

	w.NodeQueue.Chan <- kv{Key: key, Val: value}
}

// EnqueueWay - enqueue way bytes to be saved to db
func (w *Writer) EnqueueWay(item *gosmparse.Way) {

	// encode id
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(item.ID))

	// prefix way keys with 'W' to avoid id collisions
	key := append([]byte{'W'}, idBytes...)

	// encoded path
	var value []byte

	// iterate over node refs, appending each int64 id to the value
	for _, ref := range item.NodeIDs {

		// encode id
		// @todo: use varint encoding to save bytes
		refBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(refBytes, uint64(ref))

		// append to slice
		value = append(value, refBytes...)
	}

	w.WayQueue.Chan <- kv{Key: key, Val: value}
}

// Close - close the channel and block until done
func (w *Writer) Close() {
	w.NodeQueue.Close()
	w.WayQueue.Close()
}
