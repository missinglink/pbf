package leveldb

import (
	"encoding/binary"
	"log"

	"github.com/missinglink/gosmparse"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack"
)

// WriteRelation - encode and write relation to db
func (c *Connection) WriteRelation(item gosmparse.Relation) error {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(item.ID))

	// prepend relation identifier
	key = append(prefix["relation"], key...)

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

// ReadRelation - read relation from db and decode
func (c *Connection) ReadRelation(id int64) (*gosmparse.Relation, error) {

	// encode id
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))

	// prepend relation identifier
	key = append(prefix["relation"], key...)

	// read from db
	data, err := c.DB.Get(key, nil)
	if err != nil {
		return nil, err
	}

	// decode item
	var relation gosmparse.Relation
	err = msgpack.Unmarshal(data, &relation)
	if err != nil {
		log.Println("decode failed", err)
		return nil, err
	}

	return &relation, nil
}

// IterateRelations - read all relations from db and decode one-by-one
func (c *Connection) IterateRelations(cb func(*gosmparse.Relation, error)) {

	iter := c.DB.NewIterator(util.BytesPrefix(prefix["relation"]), nil)
	for iter.Next() {

		// get key/value data
		// key := iter.Key()
		data := iter.Value()

		// decode item
		var relation gosmparse.Relation
		err := msgpack.Unmarshal(data, &relation)
		if err != nil {
			log.Println("decode failed", err)
			cb(nil, err)
			return
		}

		cb(&relation, nil)
	}

	iter.Release()
}
