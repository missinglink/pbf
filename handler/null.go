package handler

import "github.com/missinglink/gosmparse"

// Null - Null
type Null struct{}

// ReadNode - called once per node
func (n *Null) ReadNode(item gosmparse.Node) {
	// noop
}

// ReadWay - called once per way
func (n *Null) ReadWay(item gosmparse.Way) {
	// noop
}

// ReadRelation - called once per relation
func (n *Null) ReadRelation(item gosmparse.Relation) {
	// noop
}
