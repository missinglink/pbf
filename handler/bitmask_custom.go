package handler

import (
	"github.com/missinglink/pbf/lib"

	"github.com/missinglink/gosmparse"
)

// BitmaskCustom - Load all elements in to memory
type BitmaskCustom struct {
	Masks    *lib.BitmaskMap
	Features *lib.FeatureSet
}

// ReadNode - called once per node
func (b *BitmaskCustom) ReadNode(item gosmparse.Node) {
	if b.Features.MatchNode(item) {
		b.Masks.Nodes.Insert(item.ID)
	}
}

// ReadWay - called once per way
func (b *BitmaskCustom) ReadWay(item gosmparse.Way) {
	if b.Features.MatchWay(item) {
		b.Masks.Ways.Insert(item.ID)

		// insert dependents in mask
		for _, ref := range item.NodeIDs {
			b.Masks.WayRefs.Insert(ref)
		}
	}
}

// ReadRelation - called once per relation
func (b *BitmaskCustom) ReadRelation(item gosmparse.Relation) {
	// @todo: relations currently not supported
	// due to requiring a 'second-pass' to gather the node ids for
	// each member way

	// if b.Features.MatchRelation(item) {
	// 	b.Masks.Relations.Insert(item.ID)
	// }
}
