package command

import (
	"log"
	"os"
	"sync"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"

	"github.com/codegangsta/cli"
)

// Nquad cli command
func Nquad(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		log.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(argv[0])

	// create parser handler
	var handle = &handler.Nquad{
		Mutex: &sync.Mutex{},
	}

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
