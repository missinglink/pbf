package handler

import (
	"strings"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/lib"
)

// Xroads - Xroads
type Xroads struct {
	Pass                 int64
	TagWhiteList         map[string]bool
	IntersectionWaysMask *lib.Bitmask
	WayNames             map[int64]string
	NodeMap              map[int64][]int64
	Coords               map[int64]*gosmparse.Node
	Mutex                *sync.Mutex
}

// ReadNode - called once per node
func (x *Xroads) ReadNode(item gosmparse.Node) {
	// skip unless this is the second pass over the file
	if x.Pass != 1 {
		return
	}

	if _, ok := x.NodeMap[item.ID]; ok {
		x.Mutex.Lock()
		x.Coords[item.ID] = &gosmparse.Node{
			Lat: item.Lat,
			Lon: item.Lon,
		}
		x.Mutex.Unlock()
	}
}

// ReadWay - called once per way
func (x *Xroads) ReadWay(item gosmparse.Way) {
	// compute intersections on first pass
	if x.Pass == 0 {

		// must be valid highway tag
		if _, ok := x.TagWhiteList[item.Tags["highway"]]; !ok {
			return
		}

		// store the way ids in an array with the nodeid as key
		for _, nodeid := range item.NodeIDs {
			x.Mutex.Lock()
			x.NodeMap[nodeid] = appendUnique(x.NodeMap[nodeid], item.ID)
			x.Mutex.Unlock()
		}
	}

	// only store names on second pass to save memory
	// (after $IntersectionWaysMask has been populated with matches)
	if x.Pass == 1 && x.IntersectionWaysMask.Has(item.ID) {
		// get the best name from the tags
		if val, ok := item.Tags["addr:street"]; ok {
			x.Mutex.Lock()
			x.WayNames[item.ID] = strings.TrimSpace(val)
			x.Mutex.Unlock()
		} else if val, ok := item.Tags["name"]; ok {
			x.Mutex.Lock()
			x.WayNames[item.ID] = strings.TrimSpace(val)
			x.Mutex.Unlock()
		} else {
			return
		}
	}
}

// ReadRelation - called once per relation
func (x *Xroads) ReadRelation(item gosmparse.Relation) {
	// noop
}

// TrimNonIntersections - remove any nodes which have less than two ways
func (x *Xroads) TrimNonIntersections() {
	x.Mutex.Lock()
	for nodeID, wayIDs := range x.NodeMap {
		if len(wayIDs) < 2 {
			delete(x.NodeMap, nodeID)
		} else {
			// insert all intersection way IDs in mask
			for _, wayID := range wayIDs {
				x.IntersectionWaysMask.Insert(wayID)
			}
		}
	}
	x.Mutex.Unlock()
}

// https://stackoverflow.com/a/9561388
func appendUnique(slice []int64, i int64) []int64 {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
