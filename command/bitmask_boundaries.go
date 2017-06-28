package command

import (
	"log"
	"os"
	"sync"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"

	"github.com/codegangsta/cli"
	"github.com/missinglink/gosmparse"
)

func recuseRelation(id int64, handle *handler.BitmaskBoundaries) {
	if members, ok := handle.RelationMembers[id]; ok {
		for _, member := range members {
			switch member.Type {
			case gosmparse.NodeType:
				handle.Masks.Nodes.Insert(member.ID)
			case gosmparse.WayType:
				handle.Masks.Ways.Insert(member.ID)
			case gosmparse.RelationType:
				if !handle.Masks.Relations.Has(member.ID) {
					handle.Masks.Relations.Insert(member.ID)
					recuseRelation(member.ID, handle)
				}
			}
		}
	} else {
		log.Println("relation not found in map", id)
	}
}

// BitmaskBoundaries cli command
func BitmaskBoundaries(c *cli.Context) error {

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// don't clobber existing bitmask file
	if _, err := os.Stat(c.Args()[1]); err == nil {
		log.Println("bitmask file already exists; don't want to override it")
		os.Exit(1)
	}

	// open database for writing
	handle := &handler.BitmaskBoundaries{
		Pass:            0,
		Mutex:           &sync.Mutex{},
		Masks:           lib.NewBitmaskMap(),
		RelationMembers: make(map[int64][]gosmparse.RelationMember),
	}

	// write to disk
	defer handle.Masks.WriteToFile(c.Args()[1])

	// Parse will block until it is done or an error occurs.
	parser.Parse(handle)

	// recurse super-relations
	for id := range handle.RelationMembers {
		if handle.Masks.Relations.Has(id) {
			recuseRelation(id, handle)
		}
	}

	// reset and add all nodes for ways in bitmask
	parser.Reset()
	handle.Pass = 1
	parser.Parse(handle)

	return nil
}
