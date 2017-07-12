package spatialite

import (
	"database/sql"

	_ "github.com/shaxbee/go-spatialite" // required database driver
)

// Connection - Connection
type Connection struct {
	DB *sql.DB
}

// Open - open connection and set up
func (c *Connection) Open(path string) {
	db, err := sql.Open("spatialite", path)
	if err != nil {
		panic(err)
	}
	c.DB = db

	// // init spatial metadata
	// _, err = db.Exec("SELECT InitSpatialMetadata(1)")
	// if err != nil {
	// 	panic(err)
	// }
}

// Close - close connection and clean up
func (c *Connection) Close() {
	defer c.DB.Close()
}
