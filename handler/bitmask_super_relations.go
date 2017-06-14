package handler

import (
	"github.com/missinglink/pbf/lib"

	"github.com/missinglink/gosmparse"
)

// BitmaskSuperRelations - Load all elements in to memory
type BitmaskSuperRelations struct {
	Masks *lib.BitmaskMap
}

// ReadNode - called once per node
func (b *BitmaskSuperRelations) ReadNode(item gosmparse.Node) { /* noop */ }

// ReadWay - called once per way
func (b *BitmaskSuperRelations) ReadWay(item gosmparse.Way) { /* noop */ }

// ReadRelation - called once per relation
func (b *BitmaskSuperRelations) ReadRelation(item gosmparse.Relation) {

	// if super relation (contains at least one other relation), add to mask
	for _, member := range item.Members {
		switch member.Type {
		case gosmparse.RelationType:
			b.Masks.Relations.Insert(item.ID)
			break
		}
	}
}
