package handler

import (
	"sync"

	"github.com/thomersch/gosmparse"
)

// ReadAll - Load all elements in to memory
type ReadAll struct {
	Mutex     *sync.Mutex
	Nodes     map[int64]gosmparse.Node
	Ways      map[int64]gosmparse.Way
	Relations map[int64]gosmparse.Relation
}

// ReadNode - called once per node
func (n *ReadAll) ReadNode(item gosmparse.Node) {
	n.Mutex.Lock()
	n.Nodes[item.ID] = item
	n.Mutex.Unlock()
}

// ReadWay - called once per way
func (n *ReadAll) ReadWay(item gosmparse.Way) {
	n.Mutex.Lock()
	n.Ways[item.ID] = item
	n.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (n *ReadAll) ReadRelation(item gosmparse.Relation) {
	n.Mutex.Lock()
	n.Relations[item.ID] = item
	n.Mutex.Unlock()
}
