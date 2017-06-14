package proxy

import (
	"github.com/missinglink/pbf/lib"

	"github.com/missinglink/gosmparse"
)

// WhiteList - filter only elements that appear in masks
type WhiteList struct {
	Handler      gosmparse.OSMReader
	NodeMask     *lib.Bitmask
	WayMask      *lib.Bitmask
	RelationMask *lib.Bitmask
}

// ReadNode - called once per node
func (p *WhiteList) ReadNode(item gosmparse.Node) {
	if nil != p.NodeMask && p.NodeMask.Has(item.ID) {
		p.Handler.ReadNode(item)
	}
}

// ReadWay - called once per way
func (p *WhiteList) ReadWay(item gosmparse.Way) {
	if nil != p.WayMask && p.WayMask.Has(item.ID) {
		p.Handler.ReadWay(item)
	}
}

// ReadRelation - called once per relation
func (p *WhiteList) ReadRelation(item gosmparse.Relation) {
	if nil != p.RelationMask && p.RelationMask.Has(item.ID) {
		p.Handler.ReadRelation(item)
	}
}
