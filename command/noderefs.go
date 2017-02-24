package command

import (
	"fmt"
	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/parser"
	"sort"

	"github.com/codegangsta/cli"
)

// NodeRefs cli command
func NodeRefs(c *cli.Context) error {

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// stats handler
	stats := &handler.Refs{Counts: make(map[int64]int)}

	// parse file
	parser.Parse(stats)

	cc := make(map[int]int)

	for _, c := range stats.Counts {
		// fmt.Printf("%v: %v\n", nodeid, c)
		cc[c]++
	}

	// sort in a slice
	var keys []int
	for k := range cc {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Println("Amount of ways that share a nodeid.")
	fmt.Println("left: number of times *any* unique node id appeared in *any* way")
	fmt.Println("right: total time the above condition occured in the file")
	for _, k := range keys {
		fmt.Printf("%3v: %v\n", k, cc[k])
	}

	return nil
}
