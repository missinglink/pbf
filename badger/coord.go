package badger

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/dgraph-io/badger/badger"
	"github.com/missinglink/gosmparse"
)

// WriteCoord - encode and write lat/lon pair to db
func (c *Connection) WriteCoord(item gosmparse.Node) error {

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

	// write to db
	err := c.DB.Set(key, value)
	log.Println("wrote", item.ID)
	if err != nil {
		log.Println("write error", err)
		return err
	}

	return nil
}

// ReadCoord - read lat/lon pair from db and decode
func (c *Connection) ReadCoord(id int64) (*gosmparse.Node, error) {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))

	// read from db
	var item badger.KVItem
	err := c.DB.Get(key, &item)
	if err != nil {
		return nil, err
	}

	data := item.Value()

	// decode item
	return &gosmparse.Node{
		ID:  id,
		Lat: math.Float64frombits(binary.BigEndian.Uint64(data[:8])),
		Lon: math.Float64frombits(binary.BigEndian.Uint64(data[8:])),
	}, nil
}
