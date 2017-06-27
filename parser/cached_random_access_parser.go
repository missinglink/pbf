package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/lib"
)

// CachedRandomAccessParser - struct to handle random access lookups to a pbf
type CachedRandomAccessParser struct {
	Parser
	Index   *gosmparse.BlobIndex
	Mutex   *sync.Mutex
	Cache   *CoordCache
	Handler *CoordCacheHandler
}

func getCacheSize() int {

	// load cache size from ENV variable
	if ff := os.Getenv("CACHE_SIZE"); ff != "" {
		i, err := strconv.Atoi(ff)
		if nil == err {
			log.Printf("custom cache size: %d\n", i)
			return i
		}
	}

	// return default size
	// tested locally a setting of 5M used ~3.3GB RAM
	return 5000000
}

// NewCachedRandomAccessParser -
func NewCachedRandomAccessParser(path string, idxPath string) *CachedRandomAccessParser {

	// load index
	index := &gosmparse.BlobIndex{}
	index.ReadFromFile(idxPath)

	// init cache
	var cache = &CoordCache{
		Mutex:      &sync.Mutex{},
		Size:       getCacheSize(),
		ClearRatio: 0.8,
		Coords:     make(map[int64]*gosmparse.Node),
		SeenMask:   lib.NewBitMask(),
	}

	var p = &CachedRandomAccessParser{
		Index: index,
		Mutex: &sync.Mutex{},
		Cache: cache,
		Handler: &CoordCacheHandler{
			Cache: cache,
		},
	}

	p.open(path)
	return p
}

// ReadNode - fetch a single node
func (p *CachedRandomAccessParser) ReadNode(osmID int64) (*gosmparse.Node, error) {

	p.Cache.Mutex.Lock()
	defer p.Cache.Mutex.Unlock()

	// log.Println("ReadNode", osmID)

	// check if we have this element in the cache
	coord, found := p.Cache.Get(osmID)
	if found {
		return coord, nil
	}

	// else load from file
	p.loadBlob("node", osmID)

	// check if we have this element in the cache now!
	coord, found = p.Cache.Get(osmID)
	if found {
		return coord, nil
	}

	// log.Printf("node not found on second pass: %d\n", osmID)
	return coord, fmt.Errorf("node not found: %d", osmID)
}

// loadBlob - fetch blob and cache returned elements
func (p *CachedRandomAccessParser) loadBlob(osmType string, osmID int64) error {

	// find the location of this element in file
	offsets, err := p.Index.BlobOffsets(osmType, osmID)
	if nil != err {
		log.Printf("target element: %s %d not found in file\n", osmType, osmID)
		return err
	}

	// log.Println("offsets", offsets)

	// Parse will block until it is done or an error occurs.
	p.Mutex.Lock()
	for _, offset := range offsets {
		p.ParseBlob(p.Handler, offset)
	}
	p.Mutex.Unlock()

	return nil
}
