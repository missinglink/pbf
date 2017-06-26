package command

import (
	"fmt"
	"time"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/parser"

	"github.com/codegangsta/cli"
)

// Stats cli command
func Stats(c *cli.Context) error {

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// stats handler
	stats := &handler.Stats{}

	// live stats
	if c.Int("interval") > 0 {
		go func() {
			for range time.Tick(time.Duration(c.Int("interval")) * time.Millisecond) {
				stats.Print()
				fmt.Println()
			}
		}()
	}

	// Parse will block until it is done or an error occurs.
	parser.Parse(stats)

	// print final stats
	stats.Print()

	// print final stats
	for _, info := range parser.GetDecoder().Index.Blobs {

		fmt.Printf("start: %v, size: %v\n", info.Start, info.Size)

		for _, group := range info.Groups {
			fmt.Printf("  type: %v, count: %v, low: %v, high: %v\n", group.Type, group.Count, group.Low, group.High)
		}
	}

	return nil
}
