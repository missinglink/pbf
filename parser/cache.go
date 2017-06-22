package parser

import (
	"log"
	"math"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/lib"
)

// CoordCache - in-memory element cache
type CoordCache struct {
	Mutex          *sync.Mutex
	Size           int
	ClearRatio     float64
	Coords         map[int64]*gosmparse.Node
	Fifo           []int64
	SeenMask       *lib.Bitmask
	DuplicatesMask *lib.Bitmask
}

// Set - store a single record in the cache
func (c *CoordCache) Set(id int64, item gosmparse.Node) {
	// c.Mutex.Lock()
	// defer c.Mutex.Unlock()

	// element already exists in cache
	if _, ok := c.Coords[id]; ok {
		return
	}

	// append id to first-in-first-out queue
	// fmt.Printf("set %d %f %f\n", id, item.Lat, item.Lon)
	c.Fifo = append(c.Fifo, id)

	// cache is full
	if len(c.Fifo) > c.Size {
		log.Println("cache purge")
		log.Println("cache size", len(c.Coords))

		var toDelete []int64                                                              // slice of ids we are going to delete this GC cycle
		var totalEntriesToDelete = int(math.Ceil(float64(c.Size) * (1.0 - c.ClearRatio))) // total entries we would like to remove in this pass

		// first purge entries we have already seen and not in the duplicate mask
		if nil != c.DuplicatesMask {
			for _, checkID := range c.Fifo {
				// log.Println(checkID, c.SeenMask.Has(checkID), !c.DuplicatesMask.Has(checkID))
				if c.SeenMask.Has(checkID) && !c.DuplicatesMask.Has(checkID) {
					toDelete = append(toDelete, checkID)
				}
			}
			log.Println("singletons deleted", len(toDelete))
		}

		// next purge oldest records first (if required)
		lastIndex := totalEntriesToDelete - len(toDelete)
		if lastIndex > 0 {
			toDelete = append(toDelete, c.Fifo[0:lastIndex]...)
		}

		log.Println("total deleted", len(toDelete))

		// perform the deletions
		c.Fifo = difference(c.Fifo, toDelete)
		for _, deadID := range toDelete {
			delete(c.Coords, deadID)
		}
		log.Println("cache size", len(c.Coords))
	}

	// set map key
	c.Coords[id] = &item
}

// Get - fetch a single record from the cache
func (c *CoordCache) Get(id int64) (*gosmparse.Node, bool) {
	// log.Println("get", id)

	coord, ok := c.Coords[id]

	if ok {
		c.SeenMask.Insert(id)
	}

	return coord, ok
}

// --- handler ---

// CoordCacheHandler - CoordCacheHandler
type CoordCacheHandler struct {
	Cache *CoordCache
	Mask  *lib.Bitmask
}

// ReadNode - called once per node
func (h *CoordCacheHandler) ReadNode(item gosmparse.Node) {

	// if a mask was supplied, use it
	if nil != h.Mask && !h.Mask.Has(item.ID) {
		return
	}

	h.Cache.Set(item.ID, gosmparse.Node{
		Lat: item.Lat,
		Lon: item.Lon,
	})
}

// ReadWay - called once per way
func (h *CoordCacheHandler) ReadWay(item gosmparse.Way) {
	// noop
}

// ReadRelation - called once per relation
func (h *CoordCacheHandler) ReadRelation(item gosmparse.Relation) {
	// noop
}

// set difference -- return all elements in A that are not present in B
// note: also ensures uniqueness
func difference(A []int64, B []int64) []int64 {
	var seen = make(map[int64]bool)
	var uniq = make(map[int64]bool)
	for _, bb := range B {
		seen[bb] = true
	}
	var C []int64
	for _, aa := range A {
		// enforce a->b exclusivity
		if _, ok := seen[aa]; !ok {
			// enforce uniqueness
			if _, ok2 := uniq[aa]; !ok2 {
				uniq[aa] = true
				C = append(C, aa)
			}
		}
	}
	return C
}
