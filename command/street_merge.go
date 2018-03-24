package command

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/missinglink/pbf/sqlite"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"
	"github.com/missinglink/pbf/tags"

	"github.com/codegangsta/cli"
	geo "github.com/paulmach/go.geo"
)

type street struct {
	Path *geo.Path
	Name string
}

type config struct {
	Format          string
	Delim           string
	ExtendedColumns bool
}

func (s *street) Print(conf *config) {

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

	var cols []string

	switch conf.Format {
	case "geojson":
		bytes, err := s.Path.ToGeoJSON().MarshalJSON()
		if nil != err {
			log.Println("failed to marshal geojson")
			os.Exit(1)
		}
		cols = append(cols, string(bytes))
	case "wkt":
		cols = append(cols, s.Path.ToWKT())
	default:
		cols = append(cols, s.Path.Encode(1.0e6))
	}

	// extended columns
	if true == conf.ExtendedColumns {
		// mid-point centroid
		var centroid = s.Path.Interpolate(0.5)
		cols = append(cols, strconv.FormatFloat(centroid.Lng(), 'f', 7, 64))
		cols = append(cols, strconv.FormatFloat(centroid.Lat(), 'f', 7, 64))

		// geodesic distance in meters
		cols = append(cols, strconv.FormatFloat(s.Path.GeoDistance(), 'f', 0, 64))

		// bounds
		var bounds = s.Path.Bound()
		var sw = bounds.SouthWest()
		var ne = bounds.NorthEast()
		cols = append(cols, strconv.FormatFloat(sw.Lng(), 'f', 7, 64))
		cols = append(cols, strconv.FormatFloat(sw.Lat(), 'f', 7, 64))
		cols = append(cols, strconv.FormatFloat(ne.Lng(), 'f', 7, 64))
		cols = append(cols, strconv.FormatFloat(ne.Lat(), 'f', 7, 64))
	}

	cols = append(cols, s.Name)
	fmt.Println(strings.Join(cols, conf.Delim))
}

// StreetMerge cli command
func StreetMerge(c *cli.Context) error {
	// config
	var conf = &config{
		Format:          "polyline",
		Delim:           "\x00",
		ExtendedColumns: c.Bool("extended"),
	}
	switch strings.ToLower(c.String("format")) {
	case "geojson":
		conf.Format = "geojson"
	case "wkt":
		conf.Format = "wkt"
	}
	if "" != c.String("delim") {
		conf.Delim = c.String("delim")
	}

	// open sqlite database connection
	// note: sqlite is used to store nodes and ways
	filename := lib.TempFileName("pbf_", ".temp.db")
	defer os.Remove(filename)
	conn := &sqlite.Connection{}
	conn.Open(filename)
	defer conn.Close()

	// parse
	parsePBF(c, conn)
	var streets = generateStreetsFromWays(conn)
	var joined = joinStreets(streets)

	// print streets
	for _, street := range joined {
		street.Print(conf)
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
	var distanceTolerance = 0.0005 // roughly 55 meters

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

	// output lines in consistent order
	keys := make([]string, len(nameMap))
	for k := range nameMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		var strs = nameMap[k]
		for _, str := range strs {
			if _, ok := merged[str]; !ok {
				ret = append(ret, str)
			}
		}
	}

	return ret
}

func loadStreetsFromDatabase(conn *sqlite.Connection, callback func(*sql.Rows)) {
	rows, err := conn.GetDB().Query(`
	SELECT
		ways.id,
		(
			SELECT GROUP_CONCAT(( nodes.lon || '#' || nodes.lat ))
			FROM way_nodes
			JOIN nodes ON way_nodes.node = nodes.id
			WHERE way = ways.id
			ORDER BY way_nodes.num ASC
		) AS nodeids,
		(
			SELECT value
			FROM way_tags
			WHERE ref = ways.id
			AND key = 'name'
			LIMIT 1
		) AS name
	FROM ways
	ORDER BY ways.id ASC;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		callback(rows)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func generateStreetsFromWays(conn *sqlite.Connection) []*street {
	var streets []*street

	loadStreetsFromDatabase(conn, func(rows *sql.Rows) {

		var wayid int
		var nodeids, name string
		err := rows.Scan(&wayid, &nodeids, &name)
		if err != nil {
			log.Fatal(err)
		}

		var wayNodes = strings.Split(nodeids, ",")
		if len(wayNodes) <= 1 {
			log.Println("found 0 refs for way", wayid)
			return
		}

		var path = geo.NewPath()
		for i, node := range wayNodes {
			coords := strings.Split(node, "#")
			lon, lonErr := strconv.ParseFloat(coords[0], 64)
			lat, latErr := strconv.ParseFloat(coords[1], 64)
			if nil != lonErr || nil != latErr {
				log.Println("error parsing coordinate as float", coords)
				return
			}
			path.InsertAt(i, geo.NewPoint(lon, lat))
		}

		streets = append(streets, &street{Name: name, Path: path})
	})

	return streets
}

func parsePBF(c *cli.Context, conn *sqlite.Connection) {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		log.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// // create parser handler
	DBHandler := &handler.Sqlite3{Conn: conn}

	// create parser
	parser := parser.NewParser(c.Args()[0])

	// streets handler
	streets := &handler.Streets{
		TagWhitelist: tags.Highway(),
		NodeMask:     lib.NewBitMask(),
		DBHandler:    DBHandler,
	}

	// parse file
	parser.Parse(streets)

	// reset file
	parser.Reset()

	// create a proxy to filter elements by mask
	filterNodes := &proxy.WhiteList{
		Handler:  DBHandler,
		NodeMask: streets.NodeMask,
	}

	// parse file again
	parser.Parse(filterNodes)
}
