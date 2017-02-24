package command

import (
	"fmt"
	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/parser"
	"os"

	"github.com/codegangsta/cli"
)

// PbfIndex cli command
func PbfIndex(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 2 {
		fmt.Println("invalid arguments, expected: {pbf} {out.idx}")
		os.Exit(1)
	}

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// Parse will block until it is done or an error occurs.
	parser.Parse(&handler.Null{})

	// write out
	parser.GetDecoder().Index.WriteToFile(c.Args()[1])

	return nil
}
