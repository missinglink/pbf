package handler

import (
	"sync"

	"github.com/missinglink/gosmparse"
)

// ReadAll - Load all elements in to memory
type ReadAll struct {
	Mutex     *sync.Mutex
	DropTags  bool
	Nodes     map[int64]gosmparse.Node
	Ways      map[int64]gosmparse.Way
	Relations map[int64]gosmparse.Relation
}

// ReadNode - called once per node
func (n *ReadAll) ReadNode(item gosmparse.Node) {
	n.Mutex.Lock()
	if n.DropTags {
		item.Tags = make(map[string]string, 0) // discard tags
	}
	n.Nodes[item.ID] = item
	n.Mutex.Unlock()
}

// ReadWay - called once per way
func (n *ReadAll) ReadWay(item gosmparse.Way) {
	n.Mutex.Lock()
	if n.DropTags {
		item.Tags = make(map[string]string, 0) // discard tags
	}
	n.Ways[item.ID] = item
	n.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (n *ReadAll) ReadRelation(item gosmparse.Relation) {
	n.Mutex.Lock()
	if n.DropTags {
		item.Tags = make(map[string]string, 0) // discard tags
	}
	n.Relations[item.ID] = item
	n.Mutex.Unlock()
}
