package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"

	"github.com/codegangsta/cli"
)

// Pelias cli command
func Pelias(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		fmt.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// create parser
	p := parser.NewParser(argv[0])

	// -- bitmask --

	// bitmask is mandatory
	var bitmaskPath = c.String("bitmask")

	// bitmask file doesn't exist
	if _, err := os.Stat(bitmaskPath); err != nil {
		fmt.Println("bitmask file doesn't exist")
		os.Exit(1)
	}

	// debug
	log.Println("loaded bitmask:", bitmaskPath)

	// read bitmask from disk
	masks := lib.NewBitmaskMap()
	masks.ReadFromFile(bitmaskPath)

	// -- random access parser --

	pbfPath, _ := filepath.Abs(argv[0])
	store := parser.NewCachedRandomAccessParser(pbfPath, pbfPath+".idx")
	store.Handler.Mask = masks.WayRefs // use mask for node cache (better memory usage)

	// --

	// create parser handler
	var handle = &handler.DenormlizedJSON{
		Mutex:           &sync.Mutex{},
		Store:           store,
		ComputeCentroid: true,
		ExportLatLons:   false,
	}

	// create filter proxy
	var filter = &proxy.WhiteList{
		Handler:      handle,
		NodeMask:     masks.Nodes,
		WayMask:      masks.Ways,
		RelationMask: masks.Relations,
	}

	// Parse will block until it is done or an error occurs.
	p.Parse(filter)

	return nil
}
