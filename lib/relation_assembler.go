package lib

import (
	"bytes"
	"log"
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/leveldb"
	"sync"

	"github.com/missinglink/gosmparse"
)

// maximum number of child relations a single relation can contain
const MAX_MEMBER_RELATIONS = 65535

// RelationAssembler - struct to handle assembling relation dependencies
type RelationAssembler struct {
	Relation *gosmparse.Relation
	Conn     *leveldb.Connection
}

// GenerateJSON - generate a byte buffer of elements encoded in json, one-per-line
func (a *RelationAssembler) GenerateJSON() bytes.Buffer {

	// buffer to store the json data
	var buffer bytes.Buffer

	// keep track of sub-relations
	var current int
	var relations = make([]*gosmparse.Relation, 1)
	relations[0] = a.Relation

	// synchronize goroutines
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	// write all members and sub members to buffer
	go func() {
		defer wg.Done()
		for current < len(relations) {

			// contains too many child relations
			if len(relations) >= MAX_MEMBER_RELATIONS {
				return
			}

			writeRelation(&buffer, a, &relations, relations[current])
			current++
		}
	}()

	// done
	wg.Wait()

	// debug max sub relations for this entity
	// fmt.Printf("%d: %d\n", a.Relation.ID, len(relations))

	return buffer
}

func writeRelation(buffer *bytes.Buffer, a *RelationAssembler, relations *[]*gosmparse.Relation, item *gosmparse.Relation) {

	// write the relation json to buffer
	buffer.Write(json.RelationFromParser(*item).Bytes())
	buffer.WriteByte('\n')

	// members
	for _, mem := range item.Members {
		switch mem.Type {
		case 0:
			if node, _ := a.Conn.ReadNode(mem.ID); nil != node {

				// clear tags
				node.Tags = make(map[string]string)

				// write node json to buffer
				buffer.Write(json.NodeFromParser(*node).Bytes())
				buffer.WriteByte('\n')

			} else {
				log.Println("missing member node", mem.ID)
			}

		case 1:

			// way
			if way, _ := a.Conn.ReadWay(mem.ID); nil != way {

				// clear tags
				way.Tags = make(map[string]string)

				// write way json to buffer
				buffer.Write(json.WayFromParser(*way).Bytes())
				buffer.WriteByte('\n')

				// refs
				for _, nodeid := range way.NodeIDs {

					if node, _ := a.Conn.ReadNode(nodeid); nil != node {

						// clear tags
						node.Tags = make(map[string]string)

						// write node json to buffer
						buffer.Write(json.NodeFromParser(*node).Bytes())
						buffer.WriteByte('\n')

					} else {
						log.Println("missing ref'd node", nodeid)
					}
				}

			} else {
				log.Println("missing member way", mem.ID)
			}
		case 2:

			// super relation
			if rel, _ := a.Conn.ReadRelation(mem.ID); nil != rel {

				// skip cyclic references to parent
				if mem.Role == "subarea" {
					continue
				}

				// clear tags
				rel.Tags = make(map[string]string)

				// add child relation to queue
				*relations = append(*relations, rel)

			} else {
				log.Println("missing member relation", mem.ID)
			}
		}
	}
}
