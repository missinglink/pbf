package json

import (
	"encoding/json"
	"fmt"

	"github.com/missinglink/gosmparse"
)

// Way struct
type Way struct {
	ID   int64             `json:"id"`
	Type string            `json:"type"`
	Hash string            `json:"hash,omitempty"`
	Tags map[string]string `json:"tags,omitempty"`
	Refs []int64           `json:"nodes"`
}

// Print json
func (way Way) Print() {
	json, _ := json.Marshal(way)
	fmt.Println(string(json))
}

// Bytes - return json
func (way Way) Bytes() []byte {
	json, _ := json.Marshal(way)
	return json
}

// WayFromParser - generate a new JSON struct based off a parse struct
func WayFromParser(item gosmparse.Way) *Way {
	return &Way{
		ID:   item.ID,
		Type: "way",
		Tags: item.Tags,
		Refs: item.NodeIDs,
	}
}
