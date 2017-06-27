package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
)

// Stats cli command
func Stats(c *cli.Context) error {

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// stats handler
	// stats := &handler.Stats{}

	// live stats
	// if c.Int("interval") > 0 {
	// 	go func() {
	// 		for range time.Tick(time.Duration(c.Int("interval")) * time.Millisecond) {
	// 			stats.Print()
	// 			fmt.Println()
	// 		}
	// 	}()
	// }

	// check if a bitmask is to be used
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

	// // Parse will block until it is done or an error occurs.
	// parser.Parse(stats)
	//
	// // print final stats
	// stats.Print()

	// print final stats
	for _, info := range parser.GetDecoder().Index.Blobs {

		fmt.Printf("start: %v, size: %v\n", info.Start, info.Size)

		for _, group := range info.Groups {
			parse := false

		Loop:
			for i := group.Low; i <= group.High; i += 64 {

				switch group.Type {
				case "node":
					if has(masks.Nodes, i) || has(masks.WayRefs, i) {
						parse = true
						break Loop
					}
				case "way":
					if has(masks.Ways, i) {
						parse = true
						break Loop
					}
				case "relation":
					if has(masks.Relations, i) {
						parse = true
						break Loop
					}
				}
			}

			fmt.Printf("  type: %v, count: %v, low: %v, high: %v, parse: %t\n", group.Type, group.Count, group.Low, group.High, parse)
		}
	}

	return nil
}

func has(mask *lib.Bitmask, v int64) bool {
	if _, ok := mask.I[uint64(v)/64]; ok {
		return true
	}
	return false
}
