package command

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"
)

// StoreNodeRefs cli command
func StoreNodeRefs(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		log.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(argv[0])

	// index file is mandatory
	if nil == parser.GetDecoder().Index {
		log.Println("PBF index required, you must generate one")
		os.Exit(1)
	}

	// bitmask is mandatory
	var bitmaskPath = c.String("bitmask")
	masks := lib.NewBitmaskMap()
	masks.ReadFromFile(bitmaskPath)

	// leveldb directory is mandatory
	var leveldbPath = c.String("leveldb")
	lib.EnsureDirectoryExists(leveldbPath, "leveldb")

	// open database connection
	conn := &leveldb.Connection{}
	conn.Open(leveldbPath)
	defer conn.Close()

	// create db writer routine
	writer := leveldb.NewCoordWriter(conn)

	// ensure all node refs are written to disk before starting on the ways
	dec := parser.GetDecoder()
	dec.Triggers = []func(int, uint64){
		func(i int, offset uint64) {
			if 0 == i {
				log.Println("writer close")
				writer.Close()
				log.Println("writer closed")
			}
		},
	}

	// create store proxy
	var store = &proxy.StoreRefs{
		Handler: &handler.Null{},
		Writer:  writer,
		Masks:   masks,
	}

	// Parse will block until it is done or an error occurs.
	parser.Parse(store)

	return nil
}
