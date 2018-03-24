package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/missinglink/pbf/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "pbf"
	app.Usage = "utilities for parsing OpenStreetMap PBF files and extracting geographic data"
	app.Commands = []cli.Command{
		{
			Name:   "stats",
			Usage:  "pbf statistics",
			Flags:  []cli.Flag{cli.IntFlag{Name: "interval, i", Usage: "write stats every i milliseconds"}},
			Action: command.Stats,
		},
		{
			Name:   "json",
			Usage:  "convert to overpass json format, optionally using bitmask to filter elements",
			Flags:  []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only output element ids in bitmask"}},
			Action: command.JSON,
		},
		{
			Name:  "json-flat",
			Usage: "convert to a json format, compulsorily using bitmask to filter elements and leveldb to denormalize where possible",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "bitmask, m", Usage: "only output element ids in bitmask"},
				cli.StringFlag{Name: "leveldb, l", Usage: "location of leveldb tmp dir"},
				cli.BoolTFlag{Name: "centroid, c", Usage: "compute centroid for non-point geometries"},
				cli.BoolFlag{Name: "geohash, g", Usage: "compute geohash property for each record"},
				cli.BoolFlag{Name: "vertices, v", Usage: "also output an array of way vertices"},
			},
			Action: command.JSONFlat,
		},
		{
			Name:   "xml",
			Usage:  "convert to osm xml format, optionally using bitmask to filter elements",
			Flags:  []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only output element ids in bitmask"}},
			Action: command.XML,
		},
		{
			Name:   "opl",
			Usage:  "convert to opl, optionally using bitmask to filter elements",
			Flags:  []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only output element ids in bitmask"}},
			Action: command.OPL,
		},
		{
			Name:        "nquad",
			Usage:       "convert to nquad, optionally using bitmask to filter elements",
			Description: "the output can be imported in to dgraph using the dgraphloader utility",
			Flags:       []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only output element ids in bitmask"}},
			Action:      command.Nquad,
		},
		{
			Name:        "cypher",
			Usage:       "convert to cypher format used by the neo4j graph database, optionally using bitmask to filter elements",
			Description: "the output can be piped directly in to neo4j: `cmd | neo4j-shell -file -`",
			Flags:       []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only output element ids in bitmask"}},
			Action:      command.Cypher,
		},
		{
			Name:   "sqlite3",
			Usage:  "import elements in to sqlite3 database, optionally using bitmask to filter elements",
			Flags:  []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only import element ids in bitmask"}},
			Action: command.Sqlite3,
		},
		{
			Name:   "leveldb",
			Usage:  "import elements in to leveldb database, optionally using bitmask to filter elements",
			Flags:  []cli.Flag{cli.StringFlag{Name: "bitmask, m", Usage: "only import element ids in bitmask"}},
			Action: command.LevelDB,
		},
		{
			Name:  "genmask",
			Usage: "generate a bitmask file by specifying feature tags to match",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "config, c", Usage: "read features from config"},
				cli.BoolFlag{Name: "indexing, i", Usage: "also write PBF index file"},
			},
			Action: command.BitmaskCustom,
		},
		{
			Name:   "genmask-boundaries",
			Usage:  "generate a bitmask file containing only elements referenced by a boundary:administrative relation",
			Action: command.BitmaskBoundaries,
		},
		{
			Name:   "genmask-super-relations",
			Usage:  "generate a bitmask file containing only relations which have at least one another relation as a member",
			Action: command.BitmaskSuperRelations,
		},
		{
			Name:   "bitmask-stats",
			Usage:  "output statistics for a bitmask file",
			Action: command.BitmaskStats,
		},
		{
			Name:  "store-noderefs",
			Usage: "store all node refs in leveldb for records matching bitmask",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "bitmask, m", Usage: "only store refs in bitmask"},
				cli.StringFlag{Name: "leveldb, l", Usage: "location of leveldb tmp dir"},
			},
			Action: command.StoreNodeRefs,
		},
		{
			Name:   "boundaries",
			Usage:  "write geojson osm boundary files using a leveldb database as source",
			Action: command.BoundaryExporter,
		},
		{
			Name:   "xroads",
			Usage:  "compute street intersections",
			Action: command.Crossroads,
		},
		{
			Name:  "streets",
			Usage: "export street segments as merged linestrings, encoded in various formats",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "format, f", Usage: "select output format, one of polyline/geojson/wkt"},
				cli.StringFlag{Name: "delim, d", Usage: "change the column delimiter (default \x00)"},
				cli.BoolFlag{Name: "extended, e", Usage: "output additional columns containing centroid and distance values"},
			},
			Action: command.StreetMerge,
		},
		{
			Name:   "noderefs",
			Usage:  "count the number of times a nodeid is referenced in file",
			Action: command.NodeRefs,
		},
		{
			Name:   "index",
			Usage:  "index a pbf file and write index to disk",
			Action: command.PbfIndex,
		},
		{
			Name:   "index-info",
			Usage:  "display a visual representation of the index file",
			Action: command.PbfIndexInfo,
		},
		{
			Name:   "find",
			Usage:  "random access to pbf",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "recurse, r", Usage: "output child elements recursively"}},
			Action: command.RandomAccess,
		},
	}

	app.Run(os.Args)
}
