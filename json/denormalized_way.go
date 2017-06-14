package json

import (
	"encoding/json"
	"fmt"

	"github.com/missinglink/gosmparse"
)

// DernomalizedWay struct
type DernomalizedWay struct {
	ID       int64             `json:"id"`
	Type     string            `json:"type"`
	Tags     map[string]string `json:"tags,omitempty"`
	Centroid *LatLon           `json:"centroid"`
	LatLons  []*LatLon         `json:"nodes"`
}

// Print json
func (way DernomalizedWay) Print() {
	json, _ := json.Marshal(way)
	fmt.Println(string(json))
}

// PrintIndent json indented
func (way DernomalizedWay) PrintIndent() {
	json, _ := json.MarshalIndent(way, "", "  ")
	fmt.Println(string(json))
}

// Bytes - return json
func (way DernomalizedWay) Bytes() []byte {
	json, _ := json.Marshal(way)
	return json
}

// DernomalizedWayFromParser - generate a new JSON struct based off a parse struct
func DernomalizedWayFromParser(item gosmparse.Way) *Way {
	return &Way{
		ID:   item.ID,
		Type: "way",
		Tags: item.Tags,
		Refs: item.NodeIDs,
	}
}
