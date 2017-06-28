package badger

import "github.com/dgraph-io/badger/badger"

// Connection - Connection
type Connection struct {
	DB *badger.KV
}

// Open - open connection and set up
func (c *Connection) Open(path string) {
	opt := badger.DefaultOptions
	opt.Dir = path
	opt.ValueDir = opt.Dir
	db, err := badger.NewKV(&opt)
	if err != nil {
		panic(err)
	}
	c.DB = db
}

// Close - close connection and clean up
func (c *Connection) Close() {
	c.DB.Close()
}
