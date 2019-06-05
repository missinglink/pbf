package handler

import (
	"log"
	"reflect"
	"sync/atomic"

	"github.com/missinglink/gosmparse"
)

// Stats - Stats
type Stats struct {
	Nodes                uint64
	Ways                 uint64
	Relations            uint64
	NodesWithName        int64
	WaysWithName         int64
	RelationsWithName    int64
	NodesWithAddress     int64
	WaysWithAddress      int64
	RelationsWithAddress int64
	NodesWithNoTags      int64
	WaysWithNoTags       int64
	RelationsWithNoTags  int64
	NodesAvgTagCount     float64
	WaysAvgTagCount      float64
	RelationsAvgTagCount float64
}

// Print stats
func (s Stats) Print() {
	k := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	for i := 0; i < k.NumField(); i++ {
		log.Printf("%s: %v\n", k.Field(i).Name, v.Field(i).Interface())
	}
}

// ReadNode - called once per node
func (s *Stats) ReadNode(item gosmparse.Node) {

	// node count
	atomic.AddUint64(&s.Nodes, 1)

	// average tag count
	s.NodesAvgTagCount = ((float64(s.Nodes-1) * s.NodesAvgTagCount) + float64(len(item.Tags))) / float64(s.Nodes)

	// no tags
	if 0 == len(item.Tags) {
		s.NodesWithNoTags++
	}

	// name stats
	if _, ok := item.Tags["name"]; ok {
		s.NodesWithName++
	}

	// street stats
	if _, ok := item.Tags["addr:street"]; ok {
		if _, ok := item.Tags["addr:housenumber"]; ok {
			s.NodesWithAddress++
		}
	}
}

// ReadWay - called once per way
func (s *Stats) ReadWay(item gosmparse.Way) {

	// way count
	atomic.AddUint64(&s.Ways, 1)

	// average tag count
	s.WaysAvgTagCount = ((float64(s.Ways-1) * s.WaysAvgTagCount) + float64(len(item.Tags))) / float64(s.Ways)

	// no tags
	if 0 == len(item.Tags) {
		s.WaysWithNoTags++
	}

	if _, ok := item.Tags["name"]; ok {
		s.WaysWithName++
	}
	if _, ok := item.Tags["addr:street"]; ok {
		if _, ok := item.Tags["addr:housenumber"]; ok {
			s.WaysWithAddress++
		}
	}
}

// ReadRelation - called once per relation
func (s *Stats) ReadRelation(item gosmparse.Relation) {

	// relation count
	atomic.AddUint64(&s.Relations, 1)

	// average tag count
	s.RelationsAvgTagCount = ((float64(s.Relations-1) * s.RelationsAvgTagCount) + float64(len(item.Tags))) / float64(s.Relations)

	// no tags
	if 0 == len(item.Tags) {
		s.RelationsWithNoTags++
	}

	if _, ok := item.Tags["name"]; ok {
		s.RelationsWithName++
	}
	if _, ok := item.Tags["addr:street"]; ok {
		if _, ok := item.Tags["addr:housenumber"]; ok {
			s.RelationsWithAddress++
		}
	}
}
