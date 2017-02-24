package handler

import (
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/tags"
	"sync"

	"github.com/thomersch/gosmparse"
)

// tags that are safe to remove
var discardableTags = tags.Discardable()

// tags that are not interesting
var uninterestingTags = tags.Uninteresting()

// JSON - JSON
type JSON struct {
	Mutex *sync.Mutex
}

// ReadNode - called once per node
func (d *JSON) ReadNode(item gosmparse.Node) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// node
	json := json.Node{
		ID:   item.ID,
		Type: "node",
		Lat:  item.Lat,
		Lon:  item.Lon,
		Tags: item.Tags,
	}

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}

// ReadWay - called once per way
func (d *JSON) ReadWay(item gosmparse.Way) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// way
	json := json.Way{
		ID:   item.ID,
		Type: "way",
		Tags: item.Tags,
		Refs: item.NodeIDs,
	}

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *JSON) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// members
	members := make([]json.Member, len(item.Members))
	for i, member := range item.Members {
		members[i] = json.Member{
			ID:   member.ID,
			Type: lib.MemberType(member.Type),
			Role: member.Role,
		}
	}

	// relation
	json := json.Relation{
		ID:      item.ID,
		Type:    "relation",
		Tags:    item.Tags,
		Members: members,
	}

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}
