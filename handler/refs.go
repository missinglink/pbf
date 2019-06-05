package handler

import "github.com/missinglink/gosmparse"

// Refs - Refs
type Refs struct {
	Counts map[int64]int
}

// ReadNode - called once per node
func (s *Refs) ReadNode(item gosmparse.Node) {
	// noop
}

// ReadWay - called once per way
func (s *Refs) ReadWay(item gosmparse.Way) {

	// add node refs to map
	for _, ref := range item.NodeIDs {
		s.Counts[ref]++
	}
}

// ReadRelation - called once per relation
func (s *Refs) ReadRelation(item gosmparse.Relation) {
	// noop
}
