package handler

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/tags"

	"github.com/missinglink/gosmparse"
)

// OPL - OPL
type OPL struct {
	Mutex *sync.Mutex
}

// ReadNode - called once per node
func (d *OPL) ReadNode(item gosmparse.Node) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// node
	var parts []string

	// id
	parts = append(parts, "n"+strconv.FormatInt(item.ID, 10))

	// tags
	var tags []string
	for key, val := range item.Tags {
		tags = append(tags, key+"="+encode(val))
	}
	parts = append(parts, "T"+strings.Join(tags, ","))

	// lat/lon
	parts = append(parts, "x"+strconv.FormatFloat(float64(item.Lon), 'f', 7, 64))
	parts = append(parts, "y"+strconv.FormatFloat(float64(item.Lat), 'f', 7, 64))

	d.Mutex.Lock()
	fmt.Println(strings.Join(parts, " "))
	d.Mutex.Unlock()
}

// ReadWay - called once per way
func (d *OPL) ReadWay(item gosmparse.Way) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// way
	var parts []string

	// id
	parts = append(parts, "w"+strconv.FormatInt(item.ID, 10))

	// tags
	var tags []string
	for key, val := range item.Tags {
		tags = append(tags, key+"="+encode(val))
	}
	parts = append(parts, "T"+strings.Join(tags, ","))

	// node refs
	var refs []string
	for _, val := range item.NodeIDs {
		refs = append(refs, "n"+strconv.FormatInt(val, 10))
	}
	parts = append(parts, "N"+strings.Join(refs, ","))

	d.Mutex.Lock()
	fmt.Println(strings.Join(parts, " "))
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *OPL) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// relation
	var parts []string

	// id
	parts = append(parts, "r"+strconv.FormatInt(item.ID, 10))

	// tags
	var tags []string
	for key, val := range item.Tags {
		tags = append(tags, key+"="+encode(val))
	}
	parts = append(parts, "T"+strings.Join(tags, ","))

	// members
	var members []string
	for _, val := range item.Members {
		var prefix = string(lib.MemberType(val.Type)[0])
		var id = strconv.FormatInt(val.ID, 10)
		members = append(members, prefix+id+"@"+val.Role)
	}
	parts = append(parts, "M"+strings.Join(members, ","))

	d.Mutex.Lock()
	fmt.Println(strings.Join(parts, " "))
	d.Mutex.Unlock()
}
