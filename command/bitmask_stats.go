package command

import (
	"fmt"
	"github.com/missinglink/pbf/lib"
	"os"

	"github.com/codegangsta/cli"
)

// BitmaskStats cli command
func BitmaskStats(c *cli.Context) error {

	// file doesn't exist
	if _, err := os.Stat(c.Args()[0]); err != nil {
		fmt.Println("bitmask file doesn't exist")
		os.Exit(1)
	}

	// open mask
	m := lib.NewBitmaskMap()
	m.ReadFromFile(c.Args()[0])

	// display stats
	fmt.Printf("Nodes:     \t%d\n", m.Nodes.Len())
	fmt.Printf("Ways:      \t%d\n", m.Ways.Len())
	fmt.Printf("Relations: \t%d\n", m.Relations.Len())

	return nil
}
