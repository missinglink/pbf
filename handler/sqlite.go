package handler

import (
	"log"

	"github.com/missinglink/pbf/sqlite"
	"github.com/missinglink/pbf/tags"

	"github.com/missinglink/gosmparse"
)

// Sqlite3 - Sqlite3
type Sqlite3 struct {
	Conn *sqlite.Connection
}

// ReadNode - called once per node
func (s *Sqlite3) ReadNode(item gosmparse.Node) {

	// id, version, uid, user, timestamp, lon, lat
	_, err := s.Conn.Stmt.Node.Exec(item.ID, item.Lon, item.Lat)
	if err != nil {
		log.Println(err)
	}

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// ref, key, value
	for key, value := range item.Tags {
		_, err := s.Conn.Stmt.NodeTags.Exec(item.ID, key, value)
		if err != nil {
			log.Println(err)
		}
	}
}

// ReadWay - called once per way
func (s *Sqlite3) ReadWay(item gosmparse.Way) {

	// id, version, uid, user, timestamp
	_, err := s.Conn.Stmt.Way.Exec(item.ID)
	if err != nil {
		log.Println(err)
	}

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// ref, key, value
	for key, value := range item.Tags {
		_, err := s.Conn.Stmt.WayTags.Exec(item.ID, key, value)
		if err != nil {
			log.Println(err)
		}
	}

	// way, num, node
	for num, nodeid := range item.NodeIDs {
		_, err := s.Conn.Stmt.WayNodes.Exec(item.ID, num, nodeid)
		if err != nil {
			log.Println(err)
		}
	}
}

// ReadRelation - called once per relation
func (s *Sqlite3) ReadRelation(item gosmparse.Relation) {

	// id, version, uid, user, timestamp
	_, err := s.Conn.Stmt.Relation.Exec(item.ID)
	if err != nil {
		log.Println(err)
	}

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// ref, key, value
	for key, value := range item.Tags {
		_, err := s.Conn.Stmt.RelationTags.Exec(item.ID, key, value)
		if err != nil {
			log.Println(err)
		}
	}

	// relation, type, ref, role
	for _, member := range item.Members {
		_, err := s.Conn.Stmt.Member.Exec(item.ID, member.Type, member.ID, member.Role)
		if err != nil {
			log.Println(err)
		}
	}
}
