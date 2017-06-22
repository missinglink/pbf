package parser

import (
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/lib"
)

// CoordCache - in-memory element cache
type CoordCache struct {
	Mutex      *sync.Mutex
	Size       int
	ClearRatio float64
	Coords     map[int64]*gosmparse.Node
	Filo       []int64
}

// Set - store a single record in the cache
func (c *CoordCache) Set(id int64, item gosmparse.Node) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// element already exists in cache
	if _, ok := c.Coords[id]; ok {
		return
	}

	// append id to first-in-last-out queue
	// fmt.Printf("set %d %f %f\n", id, item.Lat, item.Lon)
	c.Filo = append(c.Filo, id)

	// cache is full
	if len(c.Filo) >= c.Size {
		var deadID int64
		// purge entries from queue to avoid out-of-memory errors
		for len(c.Filo) >= int(float64(c.Size)*c.ClearRatio) {
			deadID, c.Filo = c.Filo[0], c.Filo[1:]
			delete(c.Coords, deadID)
		}
	}

	// set map key
	c.Coords[id] = &item
}

// Get - fetch a single record from the cache
func (c *CoordCache) Get(id int64) (*gosmparse.Node, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	coord, found := c.Coords[id]
	return coord, found
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
