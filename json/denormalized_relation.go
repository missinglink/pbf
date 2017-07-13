package json

import (
	"encoding/json"
	"fmt"

	"github.com/missinglink/gosmparse"
)

// DenormalizedRelation struct
type DenormalizedRelation struct {
	ID       int64             `json:"id"`
	Type     string            `json:"type"`
	Hash     string            `json:"hash,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	Centroid *LatLon           `json:"centroid,omitempty"`
}

// Print json
func (rel DenormalizedRelation) Print() {
	json, _ := json.Marshal(rel)
	fmt.Println(string(json))
}

// PrintIndent json indented
func (rel DenormalizedRelation) PrintIndent() {
	json, _ := json.MarshalIndent(rel, "", "  ")
	fmt.Println(string(json))
}

// Bytes - return json
func (rel DenormalizedRelation) Bytes() []byte {
	json, _ := json.Marshal(rel)
	return json
}

// DenormalizedRelationFromParser - generate a new JSON struct based off a parse struct
func DenormalizedRelationFromParser(item gosmparse.Relation) *Relation {
	return &Relation{
		ID:   item.ID,
		Type: "relation",
		Tags: item.Tags,
	}
}
