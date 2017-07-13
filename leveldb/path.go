package leveldb

import (
	"encoding/binary"
	"log"

	"github.com/missinglink/gosmparse"
)

// WritePath - encode and write an array of IDs to db
func (c *Connection) WritePath(item gosmparse.Way) error {

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

	// write to db
	err := c.DB.Put(key, value, nil)
	if err != nil {
		return err
	}

	return nil
}

// ReadPath - read array of IDS from db
func (c *Connection) ReadPath(id int64) (*gosmparse.Way, error) {

	// encode id
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(id))

	// prefix way keys with 'W' to avoid id collisions
	key := append([]byte{'W'}, idBytes...)

	if id == 341691675 {
		log.Println("341691675", "load")
	}

	// read from db
	data, err := c.DB.Get(key, nil)

	if err != nil {
		return nil, err
	}

	// decode node refs
	var refs = make([]int64, 0, len(data)/8)
	for i := 0; i < len(data); i += 8 {
		refs = append(refs, int64(binary.BigEndian.Uint64(data[i:i+8])))
	}

	// decode item
	return &gosmparse.Way{
		ID:      id,
		NodeIDs: refs,
	}, nil
}
