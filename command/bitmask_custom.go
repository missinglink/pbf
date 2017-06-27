package command

import (
	"log"
	"os"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"

	"github.com/codegangsta/cli"
)

// BitmaskCustom cli command
func BitmaskCustom(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 2 {
		log.Println("invalid arguments, expected: {pbf} {mask}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// don't clobber existing bitmask file
	if _, err := os.Stat(c.Args()[1]); err == nil {
		log.Println("bitmask file already exists; don't want to override it")
		os.Exit(1)
	}

	// check config file path
	var configPath = c.String("config")
	if "" == configPath {
		log.Println("config file required, please specify one")
		os.Exit(1)
	}

	var config, configError = lib.NewFeatureSetFromJSON(configPath)
	if nil != configError {
		log.Println("config error")
		os.Exit(1)
	}

	// also perform pbf indexing
	if c.Bool("indexing") {

		// set feature flag to enable indexing code (normally turned off for performance)
		os.Setenv("INDEXING", "ON")
	}

	// open database for writing
	handle := &handler.BitmaskCustom{
		Masks:    lib.NewBitmaskMap(),
		Features: config,
	}

	// write to disk
	defer handle.Masks.WriteToFile(c.Args()[1])

	// Parse will block until it is done or an error occurs.
	parser.Parse(handle)

	return nil
}
