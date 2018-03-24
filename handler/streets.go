package handler

import (
	"github.com/missinglink/pbf/lib"

	"github.com/missinglink/gosmparse"
)

// Streets - Streets
type Streets struct {
	TagWhitelist map[string]bool
	DBHandler    *Sqlite3
	NodeMask     *lib.Bitmask
}

// ReadNode - called once per node
func (s *Streets) ReadNode(item gosmparse.Node) {
	// noop
}

// ReadWay - called once per way
func (s *Streets) ReadWay(item gosmparse.Way) {

	// must have a valid name
	if _, ok := item.Tags["name"]; !ok {
		return
	}

	// must be valid highway tag
	if _, ok := s.TagWhitelist[item.Tags["highway"]]; !ok {
		return
	}

	// remove all tags except for 'name' to conserve storage space
	item.Tags = map[string]string{"name": item.Tags["name"]}

	// add way to database
	s.DBHandler.ReadWay(item)

	// store way refs in the node mask
	for _, nodeid := range item.NodeIDs {
		s.NodeMask.Insert(nodeid)
	}
}

// ReadRelation - called once per relation
func (s *Streets) ReadRelation(item gosmparse.Relation) {
	// noop
}
