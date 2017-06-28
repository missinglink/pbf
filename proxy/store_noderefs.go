package proxy

import (
	"log"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
)

// StoreRefs - filter only elements that appear in masks
type StoreRefs struct {
	Handler gosmparse.OSMReader
	Conn    *leveldb.Connection
	Masks   *lib.BitmaskMap
}

// ReadNode - called once per node
func (s *StoreRefs) ReadNode(item gosmparse.Node) {
	if nil != s.Masks.WayRefs && s.Masks.WayRefs.Has(item.ID) {

		// write to db, removing extra fields
		err := s.Conn.WriteCoord(gosmparse.Node{
			ID:  item.ID,
			Lat: item.Lat,
			Lon: item.Lon,
		})
		if err != nil {
			log.Println(err)
		}
	}
	if nil != s.Masks.Nodes && s.Masks.Nodes.Has(item.ID) {
		s.Handler.ReadNode(item)
	}
}

// ReadWay - called once per way
func (s *StoreRefs) ReadWay(item gosmparse.Way) {
	s.Handler.ReadWay(item)
}

// ReadRelation - called once per relation
func (s *StoreRefs) ReadRelation(item gosmparse.Relation) {
	s.Handler.ReadRelation(item)
}
