package json

import (
	"encoding/json"
	"fmt"

	"github.com/missinglink/gosmparse"
)

// DenormalizedWay struct
type DenormalizedWay struct {
	ID       int64             `json:"id"`
	Type     string            `json:"type"`
	Hash     string            `json:"hash,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Centroid *LatLon           `json:"centroid,omitempty"`
	LatLons  []*LatLon         `json:"nodes,omitempty"`
}

// Print json
func (way DenormalizedWay) Print() {
	json, _ := json.Marshal(way)
	fmt.Println(string(json))
}

// PrintIndent json indented
func (way DenormalizedWay) PrintIndent() {
	json, _ := json.MarshalIndent(way, "", "  ")
	fmt.Println(string(json))
}

// Bytes - return json
func (way DenormalizedWay) Bytes() []byte {
	json, _ := json.Marshal(way)
	return json
}

// DenormalizedWayFromParser - generate a new JSON struct based off a parse struct
func DenormalizedWayFromParser(item gosmparse.Way) *Way {
	return &Way{
		ID:   item.ID,
		Type: "way",
		Tags: item.Tags,
		Refs: item.NodeIDs,
	}
}
