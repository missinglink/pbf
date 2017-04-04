package command

import (
	"fmt"
	"strings"
	"sync"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"
	"github.com/missinglink/pbf/tags"

	"github.com/codegangsta/cli"
	geo "github.com/paulmach/go.geo"
	"github.com/thomersch/gosmparse"
)

type street struct {
	Path *geo.Path
	Name string
	Ways []gosmparse.Way
}

func (s *street) Debug() {

	// geojson
	// feature := s.Path.ToGeoJSON()
	// for _, way := range s.Ways {
	// 	for k, v := range way.Tags {
	// 		feature.SetProperty(k, v)
	// 	}
	// 	feature.SetProperty("id", way.ID)
	// }
	//
	// json, _ := feature.MarshalJSON()
	// fmt.Println(string(json))

	// polyline
	var names = make(map[string]bool)
	for _, way := range s.Ways {
		for k, v := range way.Tags {
			if k == "name" {
				names[v] = true
			}
		}
	}

	var cols []string
	cols = append(cols, s.Path.Encode(1.0e6))
	for name := range names {
		cols = append(cols, name)
	}
	fmt.Printf("%s\n", strings.Join(cols, "\x00"))
}

// StreetGraph cli command
func StreetGraph(c *cli.Context) error {
	var ways, nodes = parsePBF(c)
	var streets = generateStreetsFromWays(ways, nodes)

	var joined = joinStreets(streets)

	for _, street := range joined {
		// fmt.Printf("[%s]\n\t%s\n\t%s\n\n", street.Name, street.Path.First(), street.Path.Last())
		street.Debug()
	}

	// fmt.Println(len(ways))
	// fmt.Println(len(nodes))

	return nil
}

func joinStreets(streets []*street) []*street {

	var nameMap = make(map[string][]*street)
	var ret []*street
	var merged = make(map[*street]bool)

	for _, st := range streets {
		var normName = strings.ToLower(st.Name)
		if _, ok := nameMap[normName]; !ok {
			nameMap[normName] = []*street{st}
		} else {
			nameMap[normName] = append(nameMap[normName], st)
		}
	}

	// points do not have to be exact matches
	var distanceTolerance = 0.0001

	var reversePath = func(path *geo.Path) {
		for i := path.PointSet.Length()/2 - 1; i >= 0; i-- {
			opp := path.PointSet.Length() - 1 - i
			path.PointSet[i], path.PointSet[opp] = path.PointSet[opp], path.PointSet[i]
		}
	}

	for _, strs := range nameMap {
		for i := 0; i < len(strs); i++ {
			var str1 = strs[i]

			// fmt.Println("debug", i)
			for j := 0; j < len(strs); j++ {
				var str2 = strs[j]

				if j <= i {
					continue
				}
				if _, ok := merged[str2]; ok {
					continue
				}

				if str1.Path.Last().DistanceFrom(str2.Path.First()) < distanceTolerance {

					var match = str1.Path.Last()

					// merge str2 in to str1
					for _, point := range str2.Path.PointSet {
						if point.Equals(match) {
							continue
						}
						str1.Path.Push(&point)
					}

					merged[str2] = true
					i--
					break

				} else if str1.Path.First().DistanceFrom(str2.Path.Last()) < distanceTolerance {

					var match = str1.Path.First()

					// flip str1 & str2 points
					reversePath(str1.Path)
					reversePath(str2.Path)

					// merge str2 in to str1
					for _, point := range str2.Path.PointSet {
						if point.Equals(match) {
							continue
						}
						str1.Path.Push(&point)
					}

					// flip str1 points back
					reversePath(str1.Path)
					reversePath(str2.Path)

					merged[str2] = true
					i--
					break

				} else if str1.Path.Last().DistanceFrom(str2.Path.Last()) < distanceTolerance {

					var match = str1.Path.Last()

					// flip str2 points
					reversePath(str2.Path)

					// merge str2 in to str1
					for _, point := range str2.Path.PointSet {
						if point.Equals(match) {
							continue
						}
						str1.Path.Push(&point)
					}

					// flip str2 points back
					reversePath(str2.Path)

					merged[str2] = true
					i--
					break

				} else if str1.Path.First().DistanceFrom(str2.Path.First()) < distanceTolerance {

					var match = str1.Path.First()

					// flip str1 points
					reversePath(str1.Path)

					// merge str2 in to str1
					for _, point := range str2.Path.PointSet {
						if point.Equals(match) {
							continue
						}
						str1.Path.Push(&point)
					}

					// flip str1 points back
					reversePath(str1.Path)

					merged[str2] = true
					i--
					break

				}
			}
		}
	}

	for _, strs := range nameMap {
		for _, str := range strs {
			if _, ok := merged[str]; !ok {
				ret = append(ret, str)
			}
		}
	}

	return ret
}

func generateStreetsFromWays(ways []gosmparse.Way, nodes map[int64]gosmparse.Node) []*street {
	var streets []*street

	for _, way := range ways {
		var wayNodes, _ = getRefs(way, nodes)

		if len(wayNodes) <= 1 {
			continue
		}

		var path = geo.NewPath()
		for i, node := range wayNodes {
			path.InsertAt(i, geo.NewPoint(float64(node.Lon), float64(node.Lat)))
		}
		streets = append(streets, &street{
			Name: way.Tags["name"],
			Path: path,
			Ways: []gosmparse.Way{way},
		})
	}

	return streets
}

// note: delete nodes which don't denormalize
func getRefs(way gosmparse.Way, nodes map[int64]gosmparse.Node) ([]*gosmparse.Node, error) {
	var ret []*gosmparse.Node
	for _, nodeid := range way.NodeIDs {
		// fmt.Println(reflect.TypeOf(nodeid))
		if node, ok := nodes[nodeid]; ok {
			ret = append(ret, &node)
		} else {
			fmt.Println("failed to denormalize way", way.ID, nodeid)
			// return nil, errors.New("failed to denormalize way")
		}
	}
	return ret, nil
}

func parsePBF(c *cli.Context) ([]gosmparse.Way, map[int64]gosmparse.Node) {

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// streets handler
	streets := &handler.Streets{
		TagWhitelist: tags.Highway(),
		NodeMask:     lib.NewBitMask(),
	}

	// parse file
	parser.Parse(streets)

	// reset file
	parser.Reset()

	// nodes handler
	nodes := &handler.ReadAll{
		Nodes: make(map[int64]gosmparse.Node),
		Mutex: &sync.Mutex{},
	}

	// create a proxy to filter elements by mask
	filterRels := &proxy.WhiteList{
		Handler:  nodes,
		NodeMask: streets.NodeMask,
	}

	// parse file again
	parser.Parse(filterRels)

	return streets.Ways, nodes.Nodes
}
