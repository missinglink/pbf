package command

import (
	"log"
	"os"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"
	"github.com/missinglink/pbf/sqlite"

	"github.com/codegangsta/cli"
)

// Sqlite3 cli command
func Sqlite3(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 2 {
		log.Println("invalid arguments, expected: {pbf} {sqlitedb}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(argv[0])

	// don't clobber existing db file
	if _, err := os.Stat(argv[1]); err == nil {
		log.Println("sqlite database already exists; don't want to override it")
		os.Exit(1)
	}

	// open database connection
	conn := &sqlite.Connection{}
	conn.Open(argv[1])
	defer conn.Close()

	// create parser handler
	handle := &handler.Sqlite3{Conn: conn}

	// check if a bitmask is to be used
	var bitmaskPath = c.String("bitmask")

	// not using a bitmask
	if "" == bitmaskPath {

		// Parse will block until it is done or an error occurs.
		parser.Parse(handle)

		return nil
	}

	// read bitmask from disk
	masks := lib.NewBitmaskMap()
	masks.ReadFromFile(bitmaskPath)

	// create filter proxy
	filter := &proxy.WhiteList{
		Handler:      handle,
		NodeMask:     masks.Nodes,
		WayMask:      masks.Ways,
		RelationMask: masks.Relations,
	}

	// Parse will block until it is done or an error occurs.
	parser.Parse(filter)

	return nil
}
