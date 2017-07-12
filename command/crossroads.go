package command

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/missinglink/pbf/handler"
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
		TagWhiteList:  tags.Highway(),
		WayNames:      make(map[string]string),
		InvertedIndex: make(map[string][]string),
		Mutex:         &sync.Mutex{},
	}

	// parse file
	parser.Parse(handler)

	// iterate over inverted index
	for nodeid, wayids := range handler.InvertedIndex {

		// uniqify wayids
		uniqueWayIds := uniqString(wayids)
		if len(uniqueWayIds) > 1 {

			// generate way names
			var names []string
			sort.Strings(uniqueWayIds)
			for _, wayid := range uniqueWayIds {
				names = append(names, handler.WayNames[wayid])
			}

			// only unique ones
			var uniqueNames = uniqString(names)
			sort.Strings(uniqueNames)
			if len(uniqueNames) > 1 {
				fmt.Printf("http://openstreetmap.org/node/%-15v %v\n", nodeid, strings.Join(uniqueNames, " / "))
			}
		}
	}

	return nil
}

// convenience func to uniq a set
func uniqString(list []string) []string {
	uniqueSet := make(map[string]bool)
	for _, x := range list {
		uniqueSet[x] = true
	}
	result := make([]string, 0, len(uniqueSet))
	for x := range uniqueSet {
		result = append(result, x)
	}
	return result
}
