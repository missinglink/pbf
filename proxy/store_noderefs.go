package proxy

import (
	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
)

// StoreRefs - filter only elements that appear in masks
type StoreRefs struct {
	Handler gosmparse.OSMReader
	Writer  *leveldb.CoordWriter
	Masks   *lib.BitmaskMap
}

// ReadNode - called once per node
func (s *StoreRefs) ReadNode(item gosmparse.Node) {
	if nil != s.Masks.WayRefs && s.Masks.WayRefs.Has(item.ID) {
		s.Writer.Enqueue(&item) // write to db
	}
	if nil != s.Masks.Nodes && s.Masks.Nodes.Has(item.ID) {
		s.Handler.ReadNode(item)
	}
}

// ReadWay - called once per way
func (s *StoreRefs) ReadWay(item gosmparse.Way) {
	if nil != s.Masks.Ways && s.Masks.Ways.Has(item.ID) {
		s.Handler.ReadWay(item)
	}
}

// ReadRelation - called once per relation
func (s *StoreRefs) ReadRelation(item gosmparse.Relation) {
	if nil != s.Masks.Relations && s.Masks.Relations.Has(item.ID) {
		s.Handler.ReadRelation(item)
	}
}
