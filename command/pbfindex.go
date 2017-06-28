package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/parser"

	"github.com/codegangsta/cli"
)

// PbfIndex cli command
func PbfIndex(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		log.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// set feature flag to enable indexing code (normally turned off for performance)
	os.Setenv("INDEXING", "ON")

	// create parser
	pbfPath, _ := filepath.Abs(c.Args()[0])
	parser := parser.NewParser(pbfPath)

	// Parse will block until it is done or an error occurs.
	parser.Parse(&handler.Null{})

	return nil
}

func print(typ string, str string) {
	switch typ {
	case "node": // blue
		// fmt.Printf("\033[30;48;5;4m%s\033[0m", str) // black on blue
		fmt.Printf("\033[34m%s\033[0m", str) // blue
	case "way":
		// fmt.Printf("\033[30;48;5;2m%s\033[0m", str) // black on green
		fmt.Printf("\033[32m%s\033[0m", str) // green
	case "relation":
		// fmt.Printf("\033[30;48;5;1m%s\033[0m", str) // black on red
		fmt.Printf("\033[31m%s\033[0m", str) // red
	default:
		fmt.Printf("\033[30;48;5;7m%s\033[0m", "?") // black on white
	}
}

// PbfIndexInfo cli command
func PbfIndexInfo(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		log.Println("invalid arguments, expected: {.idx file}")
		os.Exit(1)
	}

	// create parser
	idxPath, _ := filepath.Abs(c.Args()[0])

	// load index
	index := &gosmparse.BlobIndex{}
	index.ReadFromFile(idxPath)

	fmt.Println()
	var blockcounts = make(map[string]int)
	var counts = make(map[string]int)

	// print index
	for _, info := range index.Blobs {

		// fmt.Printf("▊  start: %v, size: %v\n", info.Start, info.Size)

		for _, group := range info.Groups {

			// var code = fmt.Sprintf("%d", group.Count)
			// var code = fmt.Sprintf("%d", group.Count/1000)
			var code = "⡀"
			switch group.Count / 1000 {
			case 2:
				code = "⣀"
			case 3:
				code = "⣄"
			case 4:
				code = "⣤"
			case 5:
				code = "⣦"
			case 6:
				code = "⣶"
			case 7:
				code = "⣷"
			case 8:
				code = "⣿"
			}

			// var foo = group.Count / 1000
			// code += int(foo)

			// fmt.Printf(" ")
			print(group.Type, code)
			// fmt.Printf("\n%v\n", fmt.Sprintf("\\u%X", code))
			// fmt.Printf("\n%s\n", "\u2581")
			blockcounts[group.Type]++
			counts[group.Type] += group.Count
			fmt.Printf(" type: %v, count: %v, low: %v, high: %v\n", group.Type, group.Count, group.Low, group.High)
		}
		// fmt.Printf(" ")
	}

	fmt.Println()

	fmt.Printf("\n\033[30;48;5;%dm nodes \033[0m blocks: %d, total: %d", 4, blockcounts["node"], counts["node"])
	fmt.Printf("\n\033[30;48;5;%dm ways \033[0m blocks: %d, total: %d", 2, blockcounts["way"], counts["way"])
	fmt.Printf("\n\033[30;48;5;%dm relations \033[0m blocks: %d, total: %d\n\n", 9, blockcounts["relation"], counts["relation"])

	fmt.Printf("\n\n")

	return nil
}
