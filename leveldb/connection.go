package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// key prefixes for each element type
var prefix = func() map[string][]byte {
	return map[string][]byte{
		"node":     []byte{'N'},
		"way":      []byte{'W'},
		"relation": []byte{'R'},
	}
}()

// Connection - Connection
type Connection struct {
	DB *leveldb.DB
}

// Open - open connection and set up
func (c *Connection) Open(path string) {
	db, err := leveldb.OpenFile(path, &opt.Options{
		Compression:        opt.NoCompression,
		WriteBuffer:        120 * opt.MiB,
		BlockCacheCapacity: 120 * opt.MiB,
	})
	if err != nil {
		panic(err)
	}
	c.DB = db
}

// Close - close connection and clean up
func (c *Connection) Close() {
	c.DB.Close()
}
