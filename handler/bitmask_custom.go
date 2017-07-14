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

	// nodes in feature list
	if b.Features.MatchNode(item) {
		b.Masks.Nodes.Insert(item.ID)
	}
}

// ReadWay - called once per way
func (b *BitmaskCustom) ReadWay(item gosmparse.Way) {

	// ways in feature list
	if b.Features.MatchWay(item) {

		b.Masks.Ways.Insert(item.ID)

		// insert dependents in mask
		for _, ref := range item.NodeIDs {
			b.Masks.WayRefs.Insert(ref)
		}
	}

	// ways belonging to a relation
	if b.Masks.RelWays.Has(item.ID) {

		// insert dependents in mask
		for _, ref := range item.NodeIDs {
			b.Masks.RelNodes.Insert(ref)
		}
	}
}

// ReadRelation - called once per relation
func (b *BitmaskCustom) ReadRelation(item gosmparse.Relation) {
	if b.Features.MatchRelation(item) {

		// we currently only support the 'multipolygon' type
		// see: http://wiki.openstreetmap.org/wiki/Types_of_relation
		if val, ok := item.Tags["type"]; ok && val == "multipolygon" {

			// detect relation class
			var isSuperRelation = false
			var hasNodeCentroid = false

			// iterate members once to try to classify the relation
			for _, member := range item.Members {
				switch member.Type {
				case gosmparse.RelationType:
					isSuperRelation = true
				case gosmparse.NodeType:
					switch member.Role {
					case "label":
						hasNodeCentroid = true
					case "admin_centre":
						hasNodeCentroid = true
					}
				}
			}

			// super relations are relations containing other relations
			// we currently do not support these due to their complexity
			if isSuperRelation {
				return
			}

			// iterate over relation members
			for _, member := range item.Members {

				switch member.Type {
				case gosmparse.NodeType:

					// only store nodes if they are for 'label' or 'admin_centre'
					if member.Role == "label" || member.Role == "admin_centre" {
						b.Masks.RelNodes.Insert(member.ID)
					}

				case gosmparse.WayType:

					// only store ways if we don't have the admin_centre
					if !hasNodeCentroid {
						b.Masks.RelWays.Insert(member.ID)
					}
				}
			}

			// insert relation in mask
			b.Masks.Relations.Insert(item.ID)
		}
	}
}
