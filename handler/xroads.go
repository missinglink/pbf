package handler

import (
	"strings"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/lib"
)

// Xroads - Xroads
type Xroads struct {
	Pass           int64
	TagWhiteList   map[string]bool
	WayNodesMask   *lib.Bitmask // a mask of all node ids of all matching ways
	SharedNodeMask *lib.Bitmask // a mask of only node ids parented by more than one way
	WayNames       map[int64]string
	NodeMap        map[int64][]int64
	Coords         map[int64]*gosmparse.Node
	Mutex          *sync.Mutex
}

// ReadNode - called once per node
func (x *Xroads) ReadNode(item gosmparse.Node) {
	// skip unless this is the second pass over the file
	if x.Pass != 1 {
		return
	}

	// if this node is parented by multiple ways
	// then store the coords
	if x.SharedNodeMask.Has(item.ID) {
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

	// must be valid highway tag
	if _, ok := x.TagWhiteList[item.Tags["highway"]]; !ok {
		return
	}

	// compute intersections on first pass
	if x.Pass == 0 {

		// populate two bitmasks:
		// - WayNodesMask, a mask of all node ids of all matching ways
		// - SharedNodeMask, a mask of only node ids parented by more than one way
		for _, nodeid := range item.NodeIDs {
			if x.WayNodesMask.Has(nodeid) {
				x.SharedNodeMask.Insert(nodeid)
			}

			x.WayNodesMask.Insert(nodeid)
		}
	}

	if x.Pass == 1 {
		var waySharesNodeWithAnotherWay = false

		// iterate over way nodes checking if they are shared
		// with another way
		for _, nodeid := range item.NodeIDs {
			if x.SharedNodeMask.Has(nodeid) {
				waySharesNodeWithAnotherWay = true
				// store a map of nodeid => wayid, wayid, wayid
				x.Mutex.Lock()
				x.NodeMap[nodeid] = append(x.NodeMap[nodeid], item.ID)
				x.Mutex.Unlock()
			}
		}

		if !waySharesNodeWithAnotherWay {
			return
		}

		// only store names on second pass to save memory
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
