package command

import (
	"fmt"
	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"
	"os"
	"sync"

	"github.com/codegangsta/cli"
)

// XML cli command
func XML(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		fmt.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(argv[0])

	// create parser handler
	var handle = &handler.XML{Mutex: &sync.Mutex{}}

	// check if a bitmask is to be used
	var bitmaskPath = c.String("bitmask")

	// not using a bitmask
	if "" == bitmaskPath {

		// write header
		fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
		fmt.Println("<osm version=\"0.6\" generator=\"missinglink/pbf\">")

		// Parse will block until it is done or an error occurs.
		parser.Parse(handle)

		// write footer
		fmt.Println("</osm>")

		return nil
	}

	// using a bitmask file

	// bitmask file doesn't exist
	if _, err := os.Stat(bitmaskPath); err != nil {
		fmt.Println("bitmask file doesn't exist")
		os.Exit(1)
	}

	// debug
	fmt.Println("loaded bitmask:", bitmaskPath)

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

	// write header
	fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	fmt.Println("<osm version=\"0.6\" generator=\"missinglink/pbf\">")

	// Parse will block until it is done or an error occurs.
	parser.Parse(filter)

	// write footer
	fmt.Println("</osm>")

	return nil
}
