package handler

import (
	"sync"

	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/tags"

	"github.com/missinglink/gosmparse"
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
	json := json.NodeFromParser(item)

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
	json := json.WayFromParser(item)

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *JSON) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// relation
	json := json.RelationFromParser(item)

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}
