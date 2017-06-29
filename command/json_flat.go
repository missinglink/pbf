package command

import (
	"log"
	"os"

	"github.com/missinglink/pbf/handler"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
	"github.com/missinglink/pbf/proxy"

	"github.com/codegangsta/cli"
)

// JSONFlat cli commandw
func JSONFlat(c *cli.Context) error {

	// validate args
	var argv = c.Args()
	if len(argv) != 1 {
		log.Println("invalid arguments, expected: {pbf}")
		os.Exit(1)
	}

	// create parser
	p := parser.NewParser(argv[0])

	// index file is mandatory
	if nil == p.GetDecoder().Index {
		log.Println("PBF index required, you must generate one")
		os.Exit(1)
	}

	// bitmask is mandatory
	var bitmaskPath = c.String("bitmask")
	masks := lib.NewBitmaskMap()
	masks.ReadFromFile(bitmaskPath)

	// leveldb directory is mandatory
	var leveldbPath = c.String("leveldb")
	lib.EnsureDirectoryExists(leveldbPath, "leveldb")

	// open database connection
	conn := &leveldb.Connection{}
	conn.Open(leveldbPath)
	defer conn.Close()

	// create parser handler
	var handle = &handler.DenormalizedJSON{
		Conn:            conn,
		Writer:          lib.NewBufferedWriter(),
		ComputeCentroid: c.BoolT("centroid"),
		ComputeGeohash:  c.Bool("geohash"),
		ExportLatLons:   c.Bool("vertices"),
	}

	// close the writer routine and flush
	defer handle.Writer.Close()

	// create db writer routine
	writer := leveldb.NewCoordWriter(conn)

	// ensure all node refs are written to disk before starting on the ways
	dec := p.GetDecoder()
	dec.Triggers = []func(int, uint64){
		func(i int, offset uint64) {
			if 0 == i {
				log.Println("writer close")
				writer.Close()
				log.Println("writer closed")
			}
		},
	}

	// create store proxy
	var store = &proxy.StoreRefs{
		Handler: handle,
		Writer:  writer,
		Masks:   masks,
	}

	p.Parse(store)

	// find first way offset
	// offset, err := store.Index.FirstOffsetOfType("way")
	// if nil != err {
	// 	log.Printf("target type: %s not found in file\n", "way")
	// 	os.Exit(1)
	// }

	// Parse will block until it is done or an error occurs.
	// p.ParseFrom(filter, offset)

	return nil
}
