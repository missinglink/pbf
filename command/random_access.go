package command

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/parser"

	"github.com/codegangsta/cli"
)

// RandomAccess cli command
func RandomAccess(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 3 {
		log.Println("invalid arguments, expected: {pbf} {type} {osmid}")
		os.Exit(1)
	}

	var osmtype = argv[1]
	var osmid, _ = strconv.ParseInt(argv[2], 10, 64)

	// check if we are loading children recursively
	var recurse = c.Bool("recurse")

	// random access parser
	pbfPath, _ := filepath.Abs(c.Args()[0])
	idxPath := pbfPath + ".idx"
	access := parser.NewRandomAccessParser(pbfPath, idxPath)

	var fetchNode = func(osmid int64) {
		item, err := access.GetNode(osmid)
		if nil != err {
			log.Println(err)
			return
		}
		out := json.NodeFromParser(item)
		out.Print()
	}

	var fetchWay = func(osmid int64) {
		item, err := access.GetWay(osmid)
		if nil != err {
			log.Println(err)
			return
		}

		out := json.WayFromParser(item)
		out.Print()

		if true == recurse {
			for _, nodeid := range item.NodeIDs {
				fetchNode(int64(nodeid))
			}
		}
	}

	var fetchRelation = func(osmid int64) {}
	fetchRelation = func(osmid int64) {
		item, err := access.GetRelation(osmid)
		if nil != err {
			log.Println(err)
			return
		}

		out := json.RelationFromParser(item)
		out.Print()

		if true == recurse {
			for _, member := range item.Members {
				switch member.Type {
				case 0:
					fetchNode(member.ID)
				case 1:
					fetchWay(member.ID)
				case 2:
					// skip cyclic references to parent
					if member.Role != "subarea" {
						fetchRelation(member.ID)
					}
				}
			}
		}
	}

	switch osmtype {
	case "node":
		fetchNode(osmid)
	case "way":
		fetchWay(osmid)
	case "relation":
		fetchRelation(osmid)
	default:
		log.Println("unknown member type", osmtype)
	}

	return nil
}
