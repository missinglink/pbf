package command

import (
	"database/sql"
	"fmt"
	"log"
	"math"
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

	geo "github.com/paulmach/go.geo"
	"github.com/urfave/cli"
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

// returns: nearest, distance, remainder
var nearestMatch = func(current *street, segments []*street) (*street, float64, bool, bool, []*street) {
	var chosen int                   // closest
	var distance = math.MaxFloat64   // distance from $current in meters
	var divergence = math.MaxFloat64 // divergance from previous bearing

	var prepend bool // whether to prepend the match to current
	var reverse bool // whether the chosen match should be reversed

	var head = current.Path.First()
	var headBearing = current.Path.GetAt(1).BearingTo(head)

	var tail = current.Path.Last()
	var tailBearing = current.Path.GetAt(current.Path.Length() - 2).BearingTo(tail)

	for i, segment := range segments {
		first := segment.Path.First()
		last := segment.Path.Last()

		// first->head
		if d := first.GeoDistanceFrom(head, true); d <= distance {
			b := math.Abs(head.BearingTo(first) - headBearing)
			if d < distance || b < divergence {
				distance = d
				divergence = b
				chosen = i
				reverse = true
				prepend = true
			}
		}
		// last->head
		if d := last.GeoDistanceFrom(head, true); d <= distance {
			b := math.Abs(head.BearingTo(last) - headBearing)
			if d < distance || b < divergence {
				distance = d
				divergence = b
				chosen = i
				prepend = true
			}
		}

		// tail->first
		if d := first.GeoDistanceFrom(tail, true); d <= distance {
			b := math.Abs(tail.BearingTo(first) - tailBearing)
			if d < distance || b < divergence {
				distance = d
				divergence = b
				chosen = i
			}
		}
		// tail->last
		if d := last.GeoDistanceFrom(tail, true); d <= distance {
			b := math.Abs(tail.BearingTo(last) - tailBearing)
			if d < distance || b < divergence {
				distance = d
				divergence = b
				chosen = i
				reverse = true
			}
		}
	}

	// copy pointers to new slice to avoid memory reference errors
	var remainder = make([]*street, 0, len(segments)-1)
	remainder = append(remainder, segments[:chosen]...)
	remainder = append(remainder, segments[chosen+1:]...)

	return segments[chosen], distance, reverse, prepend, remainder
}

// reverse coordinates in path
var flip = func(path *geo.Path) {
	var l = path.PointSet.Length()
	var ps = geo.NewPointSetPreallocate(l, l)
	for i := 0; i < l; i++ {
		ps.SetAt(i, path.GetAt(l-1-i))
	}
	path.PointSet = *ps
}

var flipInPlace = func(path *geo.Path) {
	sort.SliceStable(path.PointSet, func(i, j int) bool {
		return i > j
	})
}

// graftAppend $next onto the end of $current
var graftAppend = func(current *street, next *street, reverse bool) {
	var l = next.Path.Length()
	var fn = func(i int) int { return i }
	if reverse {
		fn = func(i int) int { return l - 1 - i }
	}

	for i := 0; i < l; i++ {
		var point = next.Path.GetAt(fn(i))
		if i == 0 && point.Equals(current.Path.Last()) {
			continue
		}
		current.Path.Push(point.Clone())
	}
}

// graftPrepend $next onto the start of $current
var graftPrepend = func(current *street, next *street, reverse bool) {
	var l = next.Path.Length()
	var fn = func(i int) int { return i }
	if reverse {
		fn = func(i int) int { return l - 1 - i }
	}

	var np = geo.NewPath()

	for i := 0; i < l; i++ {
		var point = next.Path.GetAt(fn(i))
		if i == (l-1) && point.Equals(current.Path.First()) {
			continue
		}
		np.Push(point.Clone())
	}

	for _, point := range current.Path.Points() {
		np.Push(point.Clone())
	}

	current.Path = np
}

// note: segments are expected to all share the same name
// but may not nessearily belong to the same linestring
var multiLineMerge = func(segments []*street, tolerance float64) (merged []*street) {
	var current, next *street
	var dist float64
	var reverse, prepend bool

	var endofworld = &street{
		Path: geo.NewPathFromXYData([][2]float64{{-180, 90}, {-180, 90}}),
	}

	for len(segments) > 0 {

		// select starting point
		current, _, reverse, _, segments = nearestMatch(endofworld, segments)
		if reverse {
			flipInPlace(current.Path)
		}

		// select next closest segment
		for len(segments) > 0 {
			next, dist, reverse, prepend, segments = nearestMatch(current, segments)

			// nothing within tolerance distance
			if dist > tolerance {
				segments = append(segments, next) // $next wasn't matched
				break
			}

			if !prepend {
				graftAppend(current, next, reverse)
			} else {
				graftPrepend(current, next, reverse)
			}
		}

		merged = append(merged, current)
	}

	return
}

func joinStreets(streets []*street) []*street {

	var nameMap = make(map[string][]*street)
	var ret []*street
	var merged = make(map[*street]bool)

	for _, st := range streets {

		// normalize street names
		var normName = strings.ToLower(st.Name)
		if _, ok := nameMap[normName]; !ok {
			nameMap[normName] = []*street{st}
		} else {
			nameMap[normName] = append(nameMap[normName], st)
		}
	}

	// points do not have to be exact matches
	var distanceTolerance = 3.65 * 6 // width of 6 lanes (in meters)

	for norm, strs := range nameMap {
		if len(strs) > 1 {
			nameMap[norm] = multiLineMerge(strs, distanceTolerance)
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
		var maybeNodeIds sql.NullString
		err := rows.Scan(&wayid, &maybeNodeIds, &name)
		if err != nil {
			log.Fatal(err)
		}

		// handle the case where nodeids is NULL
		// note: this can occur when another tool has stripped the
		// nodes but left the ways which reference them in the file.
		if !maybeNodeIds.Valid {
			log.Println("invalid way, nodes not included in file", wayid)
			return
		}

		// convert sql.NullString to string
		if val, err := maybeNodeIds.Value(); err == nil {
			nodeids = val.(string)
		} else {
			log.Fatal("invalid nodeid value", wayid)
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
