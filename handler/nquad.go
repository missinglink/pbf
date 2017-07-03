package handler

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/tags"
)

// Nquad - Nquad
type Nquad struct {
	Mutex *sync.Mutex
}

// ReadNode - called once per node
func (d *Nquad) ReadNode(item gosmparse.Node) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	var buffer bytes.Buffer

	// tags
	for key, value := range item.Tags {
		value = strings.Replace(value, "\\", "", -1)
		value = strings.Replace(value, `"`, `\"`, -1)
		value = strings.Replace(value, "\r\n", " ", -1)
		value = strings.Replace(value, "\n", " ", -1)
		fmt.Fprintf(&buffer, "<node.%d> <%s> \"%s\" .\n", item.ID, key, value)
	}

	// flush to stdout
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}

// ReadWay - called once per way
func (d *Nquad) ReadWay(item gosmparse.Way) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	var buffer bytes.Buffer

	// tags
	for key, value := range item.Tags {
		value = strings.Replace(value, "\\", "", -1)
		value = strings.Replace(value, `"`, `\"`, -1)
		value = strings.Replace(value, "\r\n", " ", -1)
		value = strings.Replace(value, "\n", " ", -1)
		fmt.Fprintf(&buffer, "<way.%d> <%s> \"%s\" .\n", item.ID, key, value)
	}

	// refs
	for _, ref := range item.NodeIDs {
		fmt.Fprintf(&buffer, "<way.%d> <_ref> <node:%d> .\n", item.ID, ref)
	}

	// flush to stdout
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *Nquad) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	var buffer bytes.Buffer

	// tags
	for key, value := range item.Tags {
		value = strings.Replace(value, "\\", "", -1)
		value = strings.Replace(value, `"`, `\"`, -1)
		value = strings.Replace(value, "\r\n", " ", -1)
		value = strings.Replace(value, "\n", " ", -1)
		fmt.Fprintf(&buffer, "<relation.%d> <%s> \"%s\" .\n", item.ID, key, value)
	}

	// members
	for _, member := range item.Members {

		// remove difficult chars
		var role = member.Role
		role = strings.Replace(role, "'", "", -1)
		role = strings.Replace(role, `"`, "", -1)
		role = strings.Replace(role, "\\", "", -1)

		switch member.Type {
		case gosmparse.NodeType:
			fmt.Fprintf(&buffer, "<relation.%d> <_%s> <node.%d> .\n", item.ID, role, member.ID)
		case gosmparse.WayType:
			fmt.Fprintf(&buffer, "<relation.%d> <_%s> <way.%d> .\n", item.ID, role, member.ID)
		case gosmparse.RelationType:
			fmt.Fprintf(&buffer, "<relation.%d> <_%s> <relation.%d> .\n", item.ID, role, member.ID)
		}
	}

	// flush to stdout
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}
