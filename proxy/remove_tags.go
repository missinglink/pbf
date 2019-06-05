package proxy

import "github.com/missinglink/gosmparse"

// RemoveTags - remove all tags from certain element types
type RemoveTags struct {
	Handler   gosmparse.OSMReader
	Nodes     bool
	Ways      bool
	Relations bool
}

// ReadNode - called once per node
func (p *RemoveTags) ReadNode(item gosmparse.Node) {
	if true == p.Nodes {
		item.Tags = make(map[string]string)
	}
	p.Handler.ReadNode(item)
}

// ReadWay - called once per way
func (p *RemoveTags) ReadWay(item gosmparse.Way) {
	if true == p.Ways {
		item.Tags = make(map[string]string)
	}
	p.Handler.ReadWay(item)
}

// ReadRelation - called once per relation
func (p *RemoveTags) ReadRelation(item gosmparse.Relation) {
	if true == p.Relations {
		item.Tags = make(map[string]string)
	}
	p.Handler.ReadRelation(item)
}
