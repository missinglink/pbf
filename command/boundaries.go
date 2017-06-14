package command

//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"github.com/missinglink/pbf/handler"
// 	"github.com/missinglink/pbf/lib"
// 	"github.com/missinglink/pbf/proxy"
//
// 	"github.com/codegangsta/cli"
// 	geojson "github.com/paulmach/go.geojson"
// 	"github.com/missinglink/gosmparse"
// )
//
// // custom handler
// type masks struct {
// 	NodeMask     *lib.Bitmask
// 	WayMask      *lib.Bitmask
// 	RelationMask *lib.Bitmask
// 	BoundaryMask *lib.Bitmask
// }
//
// // ReadNode - called once per node
// func (m *masks) ReadNode(item gosmparse.Node) { /* noop */ }
//
// // ReadWay - called once per way
// func (m *masks) ReadWay(item gosmparse.Way) {
//
// 	// must be in bitmask
// 	if !m.WayMask.Has(item.ID) {
// 		return
// 	}
//
// 	for _, ref := range item.NodeIDs {
// 		m.NodeMask.Insert(ref)
// 	}
// }
//
// // ReadRelation - called once per relation
// func (m *masks) ReadRelation(item gosmparse.Relation) {
//
// 	// must be boundary:administrative
// 	if val, ok := item.Tags["boundary"]; !ok || "administrative" != val {
// 		return
// 	}
//
// 	// insert itself in the boundary mask
// 	m.BoundaryMask.Insert(item.ID)
//
// 	// insert it's dependents in their masks
// 	for _, member := range item.Members {
// 		switch member.Type {
// 		case 0:
// 			m.NodeMask.Insert(member.ID)
// 		case 1:
// 			m.WayMask.Insert(member.ID)
// 		case 2:
// 			m.RelationMask.Insert(member.ID)
// 		}
// 	}
// }
//
// // @todo: admin_centre tag has centroid?
// // http://www.openstreetmap.org/relation/4266321#map=11/-41.2569/174.8478
//
// // Boundaries cli command
// func Boundaries(c *cli.Context) error {
//
// 	var nMask, wMask, rMask, bMask = generateMasks(c)
//
// 	// fmt.Println(nMask.Len())
// 	// fmt.Println(wMask.Len())
// 	// fmt.Println(rMask.Len())
// 	// fmt.Println(bMask.Len())
//
// 	var nodes, ways, relations, boundaries = loadData(c, nMask, wMask, rMask, bMask)
// 	// fmt.Println(len(nodes))
// 	// fmt.Println(len(ways))
// 	// fmt.Println(len(relations))
// 	// fmt.Println(len(boundaries))
//
// 	// denormalize and output
// 	denormalize(nodes, ways, relations, boundaries)
//
// 	return nil
// }
//
// func rolesToMap(boundary gosmparse.Relation) (map[string][]int64, bool) {
//
// 	// create a map[string][]int64 of roles
// 	var roles = make(map[string][]int64)
//
// 	// iterate all members, adding them to map
// 	for _, member := range boundary.Members {
// 		if 1 != member.Type {
// 			// we only accept 'way' members at this time
// 			log.Printf("member type %d not supported\n", member.Type)
// 			continue
// 		}
// 		// todo: is this block required?
// 		if _, ok := roles[member.Role]; !ok {
// 			roles[member.Role] = make([]int64, 0)
// 		}
// 		roles[member.Role] = append(roles[member.Role], member.ID)
// 	}
//
// 	// return bool:false if no one outer role exists
// 	if _, ok := roles["outer"]; !ok {
// 		return roles, false
// 	}
//
// 	return roles, true
// }
//
// // struct to store each ring 'segment'.
// type segment struct {
// 	WayID        int64
// 	NextWayID    int64
// 	FirstNodeRef int64
// 	LastNodeRef  int64
// }
//
// // figure out orientation from two nodes
// // @todo: handle case of shared lat values
// // @todo: ring needs to be formed before this is valid?
// // func orientation(node0 gosmparse.Node, node1 gosmparse.Node) string {
// // 	var orientation = "clockwise"
// // 	if node0.Lon > 0 {
// // 		if node0.Lon > node1.Lon {
// // 			orientation = "anticlockwise"
// // 		} else if node0.Lat < node1.Lat {
// // 			orientation = "anticlockwise"
// // 		}
// // 	} else {
// // 		if node0.Lon < node1.Lon {
// // 			orientation = "anticlockwise"
// // 		} else if node0.Lat > node1.Lat {
// // 			orientation = "anticlockwise"
// // 		}
// // 	}
// // 	return orientation
// // }
//
// func mapToSegments(mmap map[string][]int64, nodes map[int64]gosmparse.Node, ways map[int64]gosmparse.Way) map[string][]segment {
//
// 	// load segments
// 	var segments = make(map[string][]segment)
// 	for role, members := range mmap {
// 		segments[role] = make([]segment, 0)
// 		for _, wayid := range members {
// 			way, ok := ways[int64(wayid)]
// 			if !ok {
// 				log.Printf("way %d not found in extract\n", wayid)
// 				continue
// 			}
// 			if len(way.NodeIDs) < 2 {
// 				log.Printf("way %d only has %d refs\n", way.ID, len(way.NodeIDs))
// 				continue
// 			}
// 			segments[role] = append(segments[role], segment{
// 				WayID:        way.ID,
// 				FirstNodeRef: way.NodeIDs[0],
// 				LastNodeRef:  way.NodeIDs[len(way.NodeIDs)-1],
// 			})
// 		}
// 	}
//
// 	var linked = lib.NewBitMask()
//
// 	// link segments
// 	for h, rs := range segments {
// 		for i, seg1 := range rs {
// 			for j, seg2 := range rs {
// 				if seg1.LastNodeRef == seg2.FirstNodeRef {
// 					if !linked.Has(seg1.LastNodeRef) {
// 						segments[h][i].NextWayID = seg2.WayID
// 						linked.Insert(seg1.LastNodeRef)
// 					}
// 				}
// 				if seg1.FirstNodeRef == seg2.LastNodeRef {
// 					if !linked.Has(seg1.FirstNodeRef) {
// 						segments[h][i].NextWayID = seg2.WayID
// 						linked.Insert(seg1.FirstNodeRef)
// 					}
// 				}
// 				if i != j { // cannot link itself on same ref
// 					if seg1.FirstNodeRef == seg2.FirstNodeRef {
// 						if !linked.Has(seg1.FirstNodeRef) {
// 							segments[h][i].NextWayID = seg2.WayID
// 							linked.Insert(seg1.FirstNodeRef)
// 						}
// 					}
// 					if seg1.LastNodeRef == seg2.LastNodeRef {
// 						if !linked.Has(seg1.LastNodeRef) {
// 							segments[h][i].NextWayID = seg2.WayID
// 							linked.Insert(seg1.LastNodeRef)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	return segments
// }
//
// func buildRings(boundary gosmparse.Relation, nodes map[int64]gosmparse.Node, ways map[int64]gosmparse.Way) [][][]float64 {
//
// 	// assemble rings
// 	var rings = make([][][]float64, 0)
//
// 	// create a map of roles
// 	var roles, ok = rolesToMap(boundary)
//
// 	// no 'outer' roles found
// 	if !ok {
// 		log.Println("no outer roles found!")
// 		return rings
// 	}
//
// 	// load segments
// 	var segments = mapToSegments(roles, nodes, ways)
//
// 	// @todo debug
// 	// var j, _ = json.MarshalIndent(segments, "", "\t")
// 	// log.Println(string(j))
//
// 	// assemble outer ring
// 	var outerNodeRefs []int64
//
// 	var role = segments["outer"]
//
// 	// some pbf extracts erroneously strip all members when doing bbox filtering
// 	if 0 == len(role) {
// 		log.Println("no roles found")
// 		return rings
// 	}
//
// 	var seg = role[0]
//
// 	var findSegmentByWayID = func(segs []segment, id int64) segment {
// 		var s segment
// 		for _, s = range segs {
// 			if s.WayID == id {
// 				return s
// 			}
// 		}
// 		return s
// 	}
//
// 	// limit max iterations to the number of members in role
// 	for x := 0; x < len(role); x++ {
//
// 		// get way data
// 		var way = ways[seg.WayID]
//
// 		// first segment
// 		var refsLen = len(outerNodeRefs)
// 		if 0 == refsLen {
// 			for _, ref := range way.NodeIDs {
// 				outerNodeRefs = append(outerNodeRefs, ref)
// 			}
// 		} else {
// 			var lastRef = outerNodeRefs[refsLen-1]
//
// 			// first node of way is the same as last node of previous way
// 			if way.NodeIDs[0] == lastRef {
// 				// skip the first; because duplicate
// 				for i := 1; i < len(way.NodeIDs); i++ {
// 					outerNodeRefs = append(outerNodeRefs, way.NodeIDs[i])
// 				}
// 				// last node of way is the same as last node of prev way (reversed)
// 			} else if way.NodeIDs[len(way.NodeIDs)-1] == lastRef {
// 				// skip the last; because duplicate
// 				for i := len(way.NodeIDs) - 2; i >= 0; i-- {
// 					outerNodeRefs = append(outerNodeRefs, way.NodeIDs[i])
// 				}
// 			}
// 		}
//
// 		// next segment in ring
// 		seg = findSegmentByWayID(role, seg.NextWayID)
// 	}
//
// 	// load coords
// 	var outer = make([][]float64, 0)
//
// 	for _, ref := range outerNodeRefs {
// 		var node = nodes[ref]
// 		outer = append(outer, []float64{float64(node.Lon), float64(node.Lat)})
// 	}
//
// 	// outer ring is first ring
// 	rings = append(rings, outer)
//
// 	// @todo assert first coord is the same as last
// 	return rings
// }
//
// func denormalize(nodes map[int64]gosmparse.Node, ways map[int64]gosmparse.Way, relations map[int64]gosmparse.Relation, boundaries map[int64]gosmparse.Relation) {
// 	for _, boundary := range boundaries {
//
// 		log.Println("boundary", boundary.ID, len(boundary.Members))
//
// 		var rings = buildRings(boundary, nodes, ways)
//
// 		var feature = geojson.NewPolygonFeature(rings)
// 		for k, v := range boundary.Tags {
// 			feature.SetProperty(k, v)
// 		}
//
// 		json, _ := json.Marshal(feature)
// 		// json, _ := json.MarshalIndent(feature, "", "\t")
// 		fmt.Println(string(json) + "\n,")
//
// 	}
// }
//
// func loadData(c *cli.Context, nodeMask *lib.Bitmask, wayMask *lib.Bitmask, relationMask *lib.Bitmask, boundaryMask *lib.Bitmask) (map[int64]gosmparse.Node, map[int64]gosmparse.Way, map[int64]gosmparse.Relation, map[int64]gosmparse.Relation) {
//
// 	log.Println("[start] loading data")
//
// 	// create parser
// 	parser := lib.NewParser(c.Args()[0])
//
// 	// relations handler
// 	rels := &handler.ReadAll{
// 		Relations: make(map[int64]gosmparse.Relation),
// 	}
//
// 	// create a proxy to filter elements by mask
// 	filterRels := &proxy.WhiteList{
// 		Handler:      rels,
// 		RelationMask: relationMask,
// 	}
//
// 	// parse file
// 	parser.Parse(filterRels)
//
// 	// reset
// 	parser.Reset()
//
// 	// elements handler
// 	elements := &handler.ReadAll{
// 		Nodes:     make(map[int64]gosmparse.Node),
// 		Ways:      make(map[int64]gosmparse.Way),
// 		Relations: make(map[int64]gosmparse.Relation),
// 	}
//
// 	// create a proxy to filter elements by mask
// 	filterAll := &proxy.WhiteList{
// 		Handler:      elements,
// 		NodeMask:     nodeMask,
// 		WayMask:      wayMask,
// 		RelationMask: boundaryMask,
// 	}
//
// 	// parse file again
// 	parser.Parse(filterAll)
// 	log.Println("[end] loading data")
//
// 	return elements.Nodes, elements.Ways, rels.Relations, elements.Relations
// }
//
// func generateMasks(c *cli.Context) (*lib.Bitmask, *lib.Bitmask, *lib.Bitmask, *lib.Bitmask) {
//
// 	// create parser
// 	parser := lib.NewParser(c.Args()[0])
//
// 	// custom handler
// 	handler := &masks{
// 		NodeMask:     lib.NewBitMask(),
// 		WayMask:      lib.NewBitMask(),
// 		RelationMask: lib.NewBitMask(),
// 		BoundaryMask: lib.NewBitMask(),
// 	}
//
// 	// parse file
// 	log.Println("[start] enumerating relations")
// 	parser.Parse(handler)
// 	log.Println("[end] enumerating relations")
//
// 	// ---------------------------------------------
//
// 	// second pass - ways only
// 	log.Println("[start] enumerating ways")
// 	parser.Reset()
// 	parser.Parse(handler)
// 	log.Println("[end] enumerating ways")
//
// 	return handler.NodeMask, handler.WayMask, handler.RelationMask, handler.BoundaryMask
// }
