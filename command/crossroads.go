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

	"github.com/mmcloughlin/geohash"
	"github.com/urfave/cli"
)

// Crossroads cli command
func Crossroads(c *cli.Context) error {
	// create parser
	parser := parser.NewParser(c.Args()[0])

	// stats handler
	handler := &handler.Xroads{
		TagWhiteList:   tags.Highway(),
		WayNodesMask:   lib.NewBitMask(),
		SharedNodeMask: lib.NewBitMask(),
		WayNames:       make(map[int64]string),
		NodeMap:        make(map[int64][]int64),
		Coords:         make(map[int64]*gosmparse.Node),
		Mutex:          &sync.Mutex{},
	}

	// parse file and compute all intersections
	parser.Parse(handler)

	// reset parser and make a second pass over the file
	// to collect the node coordinates
	parser.Reset()
	handler.Pass++
	parser.Parse(handler)

	// create a new CSV writer
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	printCSVHeader(csvWriter)

	// keep a 'global' map of all hashes seen
	// this is used to avoid duplicates
	seen := make(map[string]struct{})

	// iterate over the nodes which represent an intersection
	for nodeid, wayids := range handler.NodeMap {
		if len(wayids) > 1 {
			// write csv
			printCSVLines(csvWriter, handler, seen, nodeid, wayids)
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
func printCSVLines(csvWriter *csv.Writer, handler *handler.Xroads, seen map[string]struct{}, nodeid int64, uniqueWayIds []int64) {
	coords := handler.Coords[nodeid]

	// skip 'clipped' ways
	// ie. where the way exists in the file but not all the child nodes were included
	if coords == nil {
		log.Println("[warn] coords not found for node %d, referenced by ways: %v", nodeid, uniqueWayIds)
		return
	}

	/**
		compute a geohash of the intersection point

		note: a 6 char hash, when including its 8 neighbours offers a nice level
		of coverage which *should* hash the same for all permutations of all
		nodes at this intersection.

		note: we use neighbours rather than a common spatial prefix to avoid issues where
		the intersection spans a cell boundary and therefore the nodes dont share a common
		parent cell.

		ie. Greenwich Observatory (as an extreme worst-case example of this issue).

		see: http://missinglink.github.io/leaflet-spatial-prefix-tree/
	**/
	center := geohash.EncodeWithPrecision(coords.Lat, coords.Lon, 6)
	cells := make([]string, 0, 9)
	cells = append(cells, center)
	cells = append(cells, geohash.Neighbors(center)...)

	// generate one row per intersection
	// (there may be multiple streets intersecting a single node)
	for i, wayID1 := range uniqueWayIds {
		for j, wayID2 := range uniqueWayIds {

			// street names
			name1 := strings.TrimSpace(handler.WayNames[wayID1])
			name2 := strings.TrimSpace(handler.WayNames[wayID2])

			// normalized street names (for deduplication)
			norm1 := strings.ToLower(name1)
			norm2 := strings.ToLower(name2)

			// skip intersections of things which are the 'same'
			if j <= i || wayID1 == wayID2 || norm1 == norm2 || len(name1) == 0 || len(name2) == 0 {
				continue
			}

			// create a stable hash which can be used to deduplicate
			// multiple intersections of the same two streets
			// example of three way node: https://www.openstreetmap.org/node/26704937
			reference := []string{norm1, norm2}
			sort.Strings(reference)

			// create nine hashes which cover the center cell and its 8 neighbours
			hashes := make([]string, 0, 9)
			for _, cell := range cells {
				hashes = append(hashes, strings.Join(append(reference, cell), "_"))
			}

			// check if this intersection is a duplicate of any previously
			// computed hashes.
			isDuplicate := false
			for _, hash := range hashes {
				if _, ok := seen[hash]; ok {
					isDuplicate = true
					break
				}
			}

			// skip duplicates
			if isDuplicate {
				continue
			}

			// store hashes to avoid future duplicates.
			for _, hash := range hashes {
				seen[hash] = struct{}{}
			}

			err := csvWriter.Write([]string{
				"openstreetmap",
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
