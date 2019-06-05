package leveldb

import (
	"encoding/binary"
	"log"

	"github.com/missinglink/gosmparse"
	"github.com/vmihailenco/msgpack"
)

// WriteNode - encode and write node to db
func (c *Connection) WriteNode(item gosmparse.Node) error {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(item.ID))

	// prepend node identifier
	key = append(prefix["node"], key...)

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

// ReadNode - read node from db and decode
func (c *Connection) ReadNode(id int64) (*gosmparse.Node, error) {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))

	// prepend node identifier
	key = append(prefix["node"], key...)

	// read from db
	data, err := c.DB.Get(key, nil)
	if err != nil {
		return nil, err
	}

	// decode item
	var node gosmparse.Node
	err = msgpack.Unmarshal(data, &node)
	if err != nil {
		log.Println("decode failed", err)
		return nil, err
	}

	return &node, nil
}
