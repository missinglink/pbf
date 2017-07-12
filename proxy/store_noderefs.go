package proxy

import (
	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
)

// StoreRefs - filter only elements that appear in masks
type StoreRefs struct {
	Handler     gosmparse.OSMReader
	CoordWriter *leveldb.CoordWriter
	PathWriter  *leveldb.PathWriter
	Masks       *lib.BitmaskMap
}

// ReadNode - called once per node
func (s *StoreRefs) ReadNode(item gosmparse.Node) {
	if nil != s.Masks.WayRefs && s.Masks.WayRefs.Has(item.ID) {
		s.CoordWriter.Enqueue(&item) // write to db
	} else if nil != s.Masks.RelNodes && s.Masks.RelNodes.Has(item.ID) {
		s.CoordWriter.Enqueue(&item) // write to db
	}
	if nil != s.Masks.Nodes && s.Masks.Nodes.Has(item.ID) {
		s.Handler.ReadNode(item)
	}
}

// ReadWay - called once per way
func (s *StoreRefs) ReadWay(item gosmparse.Way) {
	if nil != s.Masks.RelWays && s.Masks.RelWays.Has(item.ID) {
		s.PathWriter.Enqueue(&item) // write to db
	}
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
