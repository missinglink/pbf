package command

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/tags"

	"github.com/codegangsta/cli"
)

// Crossroads cli command
func Crossroads(c *cli.Context) error {

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// stats handler
	handler := &handler.Xroads{
		TagWhiteList:         tags.Highway(),
		IntersectionWaysMask: lib.NewBitMask(),
		WayNames:             make(map[int64]string),
		NodeMap:              make(map[int64][]int64),
		Coords:               make(map[int64]*gosmparse.Node),
		Mutex:                &sync.Mutex{},
	}

	// parse file and compute all intersections
	parser.Parse(handler)

	// remove any nodes which are members of less than two ways
	handler.TrimNonIntersections()

	// reset parser and make a second pass over the file
	// to collect the node coordinates
	parser.Reset()
	handler.Pass++
	parser.Parse(handler)

	// create a new CSV writer
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	printCSVHeader(csvWriter)

	// iterate over the nodes which represent an intersection
	for nodeid, wayids := range handler.NodeMap {
		if len(wayids) > 1 {

			// write csv
			printCSVLines(csvWriter, handler, nodeid, wayids)
		}
	}

	return nil
}

// print the CSV header
func printCSVHeader(csvWriter *csv.Writer) {
	err := csvWriter.Write([]string{
		"source",
		"ID",
		"layer",
		"lat",
		"lon",
		"street",
		"cross_street",
	})
	if err != nil {
		fmt.Println(err)
	}
}

// print crossroad info as CSV line
func printCSVLines(csvWriter *csv.Writer, handler *handler.Xroads, nodeid int64, uniqueWayIds []int64) {
	var coords = handler.Coords[nodeid]
	var seen = make(map[string]bool)

	// generate one row per intersection
	// (there may be multiple streets intersecting a single node)
	for i, wayID1 := range uniqueWayIds {
		for j, wayID2 := range uniqueWayIds {
			var name1 = strings.TrimSpace(handler.WayNames[wayID1])
			var name2 = strings.TrimSpace(handler.WayNames[wayID2])
			if j <= i || wayID1 == wayID2 || name1 == name2 || len(name1) == 0 || len(name2) == 0 {
				continue
			}

			// create a stable identifier which can be used to deduplicate
			// multiple intersections of the same two streets
			// example of three way node: https://www.openstreetmap.org/node/26704937
			var reference = []string{name1, name2}
			sort.Strings(reference)
			var identifier = strings.Join(reference, "_")

			// skip duplicates
			if _, ok := seen[identifier]; ok {
				continue
			} else {
				seen[identifier] = true
			}

			err := csvWriter.Write([]string{
				"osm",
				fmt.Sprintf("w%d-n%d-w%d", wayID1, nodeid, wayID2),
				"intersection",
				fmt.Sprintf("%f", coords.Lat),
				fmt.Sprintf("%f", coords.Lon),
				name1,
				name2,
			})
			if err != nil {
				log.Println(err)
			}
		}
	}
}
