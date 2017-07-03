package handler

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/tags"
)

// Cypher - Cypher
type Cypher struct {
	Mutex    *sync.Mutex
	KeyRegex *regexp.Regexp
}

// ReadNode - called once per node
func (d *Cypher) ReadNode(item gosmparse.Node) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// serialize tags
	var tags []string
	for key, value := range item.Tags {
		key = d.KeyRegex.ReplaceAllString(key, "_")
		value = strings.Replace(value, "\"", "\\\"", -1)
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", key, value))
	}

	var buffer bytes.Buffer

	// node
	fmt.Fprintf(&buffer, "CREATE (N%d:Element:Node {%s});\n", item.ID, strings.Join(tags, ", "))

	// flush to stdout
	fmt.Fprintf(&buffer, "\n")
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}

// ReadWay - called once per way
func (d *Cypher) ReadWay(item gosmparse.Way) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// serialize tags
	var tags []string
	for key, value := range item.Tags {
		key = d.KeyRegex.ReplaceAllString(key, "_")
		value = strings.Replace(value, "\"", "\\\"", -1)
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", key, value))
	}

	var buffer bytes.Buffer

	// way
	fmt.Fprintf(&buffer, "CREATE (W%d:Element:Way {%s});\n", item.ID, strings.Join(tags, ", "))

	// refs
	for _, ref := range item.NodeIDs {
		fmt.Fprintf(&buffer, "CREATE (W%d)-[:CONTAINS]->(N%d:Node);\n", item.ID, ref)
	}

	// flush to stdout
	fmt.Fprintf(&buffer, "\n")
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *Cypher) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// serialize tags
	var tags []string
	for key, value := range item.Tags {
		key = d.KeyRegex.ReplaceAllString(key, "_")
		value = strings.Replace(value, "\"", "\\\"", -1)
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", key, value))
	}

	var buffer bytes.Buffer

	// relation
	fmt.Fprintf(&buffer, "CREATE (R%d:Element:Relation {%s})\n", item.ID, strings.Join(tags, ", "))

	// members
	for _, member := range item.Members {

		// remove difficult chars
		var role = member.Role
		role = strings.Replace(role, "'", "", -1)
		role = strings.Replace(role, "\\", "", -1)

		switch member.Type {
		case gosmparse.NodeType:
			// fmt.Fprintf(&buffer, "CREATE (R%d)-[:CONTAINS {role: '%s'}]->(N%d:Node);\n", item.ID, role, member.ID)
		case gosmparse.WayType:
			// fmt.Fprintf(&buffer, "CREATE (R%d)-[:CONTAINS {role: '%s'}]->(W%d:Way);\n", item.ID, role, member.ID)
		case gosmparse.RelationType:
			fmt.Fprintf(&buffer, "CREATE (R%d)-[:CONTAINS {role: '%s'}]->(R%d:Relation);\n", item.ID, role, member.ID)
		}
	}

	// flush to stdout
	fmt.Fprintf(&buffer, "\n")
	d.Mutex.Lock()
	os.Stdout.Write(buffer.Bytes())
	d.Mutex.Unlock()
}
