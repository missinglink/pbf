package handler

import (
	"strconv"
	"sync"

	"github.com/missinglink/gosmparse"
)

// Xroads - Xroads
type Xroads struct {
	TagWhiteList  map[string]bool
	WayNames      map[string]string
	InvertedIndex map[string][]string
	Mutex         *sync.Mutex
}

// ReadNode - called once per node
func (x *Xroads) ReadNode(item gosmparse.Node) {
	// noop
}

// ReadWay - called once per way
func (x *Xroads) ReadWay(item gosmparse.Way) {

	// must be valid highway tag
	if _, ok := x.TagWhiteList[item.Tags["highway"]]; !ok {
		return
	}

	// convert int64 to string
	var wayIDString = strconv.FormatInt(item.ID, 10)

	// get the best name from the tags
	if val, ok := item.Tags["addr:street"]; ok {
		x.Mutex.Lock()
		x.WayNames[wayIDString] = val
		x.Mutex.Unlock()
	} else if val, ok := item.Tags["name"]; ok {
		x.Mutex.Lock()
		x.WayNames[wayIDString] = val
		x.Mutex.Unlock()
	} else {
		return
	}

	// store the way ids in an array with the nodeid as key
	for _, nodeid := range item.NodeIDs {
		var nodeIDString = strconv.FormatInt(nodeid, 10)
		x.Mutex.Lock()
		x.InvertedIndex[nodeIDString] = append(x.InvertedIndex[nodeIDString], wayIDString)
		x.Mutex.Unlock()
	}
}

// ReadRelation - called once per relation
func (x *Xroads) ReadRelation(item gosmparse.Relation) {
	// noop
}
