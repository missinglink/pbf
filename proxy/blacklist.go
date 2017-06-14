package proxy

import (
	"github.com/missinglink/pbf/lib"

	"github.com/missinglink/gosmparse"
)

// BlackList - filter only elements that do not appear in masks
type BlackList struct {
	Handler      gosmparse.OSMReader
	NodeMask     *lib.Bitmask
	WayMask      *lib.Bitmask
	RelationMask *lib.Bitmask
}

// ReadNode - called once per node
func (p *BlackList) ReadNode(item gosmparse.Node) {
	if nil != p.NodeMask && !p.NodeMask.Has(item.ID) {
		p.Handler.ReadNode(item)
	}
}

// ReadWay - called once per way
func (p *BlackList) ReadWay(item gosmparse.Way) {
	if nil != p.WayMask && !p.WayMask.Has(item.ID) {
		p.Handler.ReadWay(item)
	}
}

// ReadRelation - called once per relation
func (p *BlackList) ReadRelation(item gosmparse.Relation) {
	if nil != p.RelationMask && !p.RelationMask.Has(item.ID) {
		p.Handler.ReadRelation(item)
	}
}
