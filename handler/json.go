package handler

import (
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/tags"

	"github.com/missinglink/gosmparse"
)

// tags that are safe to remove
var discardableTags = tags.Discardable()

// tags that are not interesting
var uninterestingTags = tags.Uninteresting()

// JSON - JSON
type JSON struct {
	Writer *lib.BufferedWriter
}

// ReadNode - called once per node
func (d *JSON) ReadNode(item gosmparse.Node) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// node
	obj := json.NodeFromParser(item)
	d.Writer.Queue <- obj.Bytes()
}

// ReadWay - called once per way
func (d *JSON) ReadWay(item gosmparse.Way) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// way
	obj := json.WayFromParser(item)
	d.Writer.Queue <- obj.Bytes()
}

// ReadRelation - called once per relation
func (d *JSON) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// relation
	obj := json.RelationFromParser(item)
	d.Writer.Queue <- obj.Bytes()
}
