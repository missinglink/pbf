package handler

import (
	"bytes"
	"fmt"
	"github.com/missinglink/pbf/lib"
	"os"
	"sync"

	"github.com/missinglink/gosmparse"
)

// XML - XML
type XML struct {
	Mutex *sync.Mutex
}

// ReadNode - called once per node
func (d *XML) ReadNode(item gosmparse.Node) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	var buffer bytes.Buffer

	// node
	fmt.Fprintf(&buffer, "\t<node id=\"%d\" lat=\"%f\" lon=\"%f\">\n", item.ID, item.Lat, item.Lon)

	// tags
	for key, val := range item.Tags {
		fmt.Fprintf(&buffer, "\t\t<tag k=\"%s\" v=\"%s\" />\n", key, val)
	}

	fmt.Fprintln(&buffer, "\t</node>")

	// flush to stdout
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}

// ReadWay - called once per way
func (d *XML) ReadWay(item gosmparse.Way) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	var buffer bytes.Buffer

	// way
	fmt.Fprintf(&buffer, "\t<way id=\"%d\">\n", item.ID)

	// refs
	for _, nodeid := range item.NodeIDs {
		fmt.Fprintf(&buffer, "\t\t<nd ref=\"%d\" />\n", nodeid)
	}

	// tags
	for key, val := range item.Tags {
		fmt.Fprintf(&buffer, "\t\t<tag k=\"%s\" v=\"%s\" />\n", key, val)
	}

	fmt.Fprintln(&buffer, "\t</way>")

	// flush to stdout
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *XML) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	var buffer bytes.Buffer

	// relation
	fmt.Fprintf(&buffer, "\t<relation id=\"%d\">\n", item.ID)

	// members
	for _, mem := range item.Members {
		fmt.Fprintf(&buffer, "\t\t<member type=\"%s\" ref=\"%d\" role=\"%s\" />\n", lib.MemberType(mem.Type), mem.ID, mem.Role)
	}

	// tags
	for key, val := range item.Tags {
		fmt.Fprintf(&buffer, "\t\t<tag k=\"%s\" v=\"%s\" />\n", key, val)
	}

	fmt.Fprintln(&buffer, "\t</relation>")

	// flush to stdout
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}
