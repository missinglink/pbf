package leveldb

import (
	"encoding/binary"
	"log"

	"github.com/missinglink/gosmparse"
	"github.com/vmihailenco/msgpack"
)

// WriteWay - encode and write way to db
func (c *Connection) WriteWay(item gosmparse.Way) error {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(item.ID))

	// prepend way identifier
	key = append(prefix["way"], key...)

	// encode item
	value, err := msgpack.Marshal(item)
	if err != nil {
		log.Println("encode failed", err)
		return err
	}

	// write to db
	err = c.DB.Put(key, value, nil)
	if err != nil {
		return err
	}

	return nil
}

// ReadWay - read way from db and decode
func (c *Connection) ReadWay(id int64) (*gosmparse.Way, error) {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))

	// prepend way identifier
	key = append(prefix["way"], key...)

	// read from db
	data, err := c.DB.Get(key, nil)
	if err != nil {
		return nil, err
	}

	// decode item
	var way gosmparse.Way
	err = msgpack.Unmarshal(data, &way)
	if err != nil {
		log.Println("decode failed", err)
		return nil, err
	}

	return &way, nil
}
