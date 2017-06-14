package lib

import "github.com/missinglink/gosmparse"

// Channels struct
type Channels struct {
	Nodes     chan gosmparse.Node
	Ways      chan gosmparse.Way
	Relations chan gosmparse.Relation
}

// Close all channels
func (c *Channels) Close() {
	close(c.Nodes)
	close(c.Ways)
	close(c.Relations)
}

// NewChannels returns a new struct of channels
// @!todo: try return as *Channels
func NewChannels() Channels {
	return Channels{
		make(chan gosmparse.Node, 256),
		make(chan gosmparse.Way, 256),
		make(chan gosmparse.Relation, 256),
	}
}

// ChannelHandler struct
type ChannelHandler struct {
	Channels Channels
}

// ReadNode - called once per node
func (h *ChannelHandler) ReadNode(n gosmparse.Node) {
	h.Channels.Nodes <- n
}

// ReadWay - called once per way
func (h *ChannelHandler) ReadWay(w gosmparse.Way) {
	h.Channels.Ways <- w
}

// ReadRelation - called once per relation
func (h *ChannelHandler) ReadRelation(r gosmparse.Relation) {
	h.Channels.Relations <- r
}
