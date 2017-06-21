package command

import (
	"fmt"
	"os"
	"path/filepath"

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
		fmt.Println("invalid arguments, expected: {pbf} {mask}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// don't clobber existing bitmask file
	if _, err := os.Stat(c.Args()[1]); err == nil {
		fmt.Println("bitmask file already exists; don't want to override it")
		os.Exit(1)
	}

	// check config file path
	var configPath = c.String("config")
	if "" == configPath {
		fmt.Println("config file required, please specify one")
		os.Exit(1)
	}

	var config, configError = lib.NewFeatureSetFromJSON(configPath)
	if nil != configError {
		fmt.Println("config error")
		os.Exit(1)
	}

	// also perform pbf indexing
	if c.Bool("indexing") {

		// set feature flag to enable indexing code (normalled turned off for performance)
		os.Setenv("INDEXING", "ON")

		pbfPath, _ := filepath.Abs(c.Args()[0])
		// write out to disk
		defer func() {
			parser.GetDecoder().Index.WriteToFile(pbfPath + ".idx")
		}()
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
