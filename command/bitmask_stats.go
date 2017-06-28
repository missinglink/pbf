package command

import (
	"log"
	"os"

	"github.com/missinglink/pbf/lib"

	"github.com/codegangsta/cli"
)

// BitmaskStats cli command
func BitmaskStats(c *cli.Context) error {

	// file doesn't exist
	if _, err := os.Stat(c.Args()[0]); err != nil {
		log.Println("bitmask file doesn't exist")
		os.Exit(1)
	}

	// open mask
	m := lib.NewBitmaskMap()
	m.ReadFromFile(c.Args()[0])

	// display stats
	m.Print()

	return nil
}
