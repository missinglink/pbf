package json

import (
	"encoding/json"
	"fmt"

	"github.com/missinglink/gosmparse"
)

// Node struct
type Node struct {
	ID   int64             `json:"id"`
	Type string            `json:"type"`
	Hash string            `json:"hash,omitempty"`
	Lat  float64           `json:"lat"`
	Lon  float64           `json:"lon"`
	Tags map[string]string `json:"tags,omitempty"`
}

// Print json
func (node Node) Print() {
	json, _ := json.Marshal(node)
	fmt.Println(string(json))
}

// Bytes - return json
func (node Node) Bytes() []byte {
	json, _ := json.Marshal(node)
	return json
}

// NodeFromParser - generate a new JSON struct based off a parse struct
func NodeFromParser(item gosmparse.Node) *Node {
	return &Node{
		ID:   item.ID,
		Type: "node",
		Lat:  truncate(item.Lat),
		Lon:  truncate(item.Lon),
		Tags: item.Tags,
	}
}
