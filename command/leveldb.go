package command

import (
	"log"
	"os"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"

	"github.com/codegangsta/cli"
)

// LevelDB cli command
func LevelDB(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 2 {
		log.Println("invalid arguments, expected: {pbf} {leveldb}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(argv[0])

	// stat leveldb destination
	lib.EnsureDirectoryExists(argv[1], "leveldb")

	// open database connection
	conn := &leveldb.Connection{}
	conn.Open(argv[1])
	defer conn.Close()

	// create parser handler
	var handle = &handler.LevelDB{Conn: conn}

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
		NodeMask:     masks.Nodes,
		WayMask:      masks.Ways,
		RelationMask: masks.Relations,
		Handler:      handle,
	}

	// Parse will block until it is done or an error occurs.
	parser.Parse(filter)

	return nil
}
