package handler

import (
	"log"

	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/tags"

	"github.com/missinglink/gosmparse"
)

// LevelDB - LevelDB
type LevelDB struct {
	Conn *leveldb.Connection
}

// ReadNode - called once per node
func (s *LevelDB) ReadNode(item gosmparse.Node) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// write to db
	err := s.Conn.WriteNode(item)
	if err != nil {
		log.Println(err)
	}
}

// ReadWay - called once per way
func (s *LevelDB) ReadWay(item gosmparse.Way) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// write to db
	err := s.Conn.WriteWay(item)
	if err != nil {
		log.Println(err)
	}
}

// ReadRelation - called once per relation
func (s *LevelDB) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// write to db
	err := s.Conn.WriteRelation(item)
	if err != nil {
		log.Println(err)
	}
}
