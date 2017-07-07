package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/missinglink/gosmparse"
	tagutils "github.com/missinglink/pbf/tags"
)

// Config - struct representing the config file format
type Config map[string]Group

// Group - a collection of patterns
type Group []Pattern

// Pattern - a collection of conditions
type Pattern []Condition

// Condition - a single match condition
type Condition string

// FeatureSet struct
type FeatureSet struct {
	NodePatterns     Group
	WayPatterns      Group
	RelationPatterns Group
}

// NewFeatureSetFromJSON - load a featureset from JSON
func NewFeatureSetFromJSON(path string) (*FeatureSet, error) {

	file, e := ioutil.ReadFile(path)
	if nil != e {
		return nil, e
	}

	decoder := json.NewDecoder(bytes.NewReader(file))
	var config Config
	decoder.Decode(&config)

	var fs = &FeatureSet{}
	fs.NodePatterns = config["node"]
	fs.WayPatterns = config["way"]
	fs.RelationPatterns = config["relation"]

	return fs, nil
}

// MatchNode - yes/no if any target feature matches this node
func (fs *FeatureSet) MatchNode(n gosmparse.Node) bool {
	return matchGroup(n.Tags, fs.NodePatterns)
}

// MatchWay - yes/no if any target feature matches this way
func (fs *FeatureSet) MatchWay(w gosmparse.Way) bool {
	return matchGroup(w.Tags, fs.WayPatterns)
}

// MatchRelation - yes/no if any target feature matches this relation
func (fs *FeatureSet) MatchRelation(r gosmparse.Relation) bool {
	return matchGroup(r.Tags, fs.RelationPatterns)
}

// matchGroup - match ANY pattern in group (logical OR)
func matchGroup(tags map[string]string, group Group) bool {

	// no tags at all
	if len(tags) == 0 {
		return false
	}

	// trim all keys/value of extra whitespace
	tags = tagutils.Trim(tags)

	// OR groups
	for _, pattern := range group {
		// AND conditions
		if matchPattern(tags, pattern) {
			return true
		}
	}

	return false
}

// matchPattern - match ALL conditions in pattern (logical AND)
func matchPattern(tags map[string]string, pattern Pattern) bool {
	if len(pattern) == 0 {
		return false
	}
	for _, condition := range pattern {
		if !matchCondition(tags, condition) {
			return false
		}
	}
	return true
}

// matchCondition - match a single condition
func matchCondition(tags map[string]string, condition Condition) bool {

	// split to array, first part will be the key and second the value (if required)
	part := strings.Split(string(condition), "=")
	val, isFound := tags[part[0]]

	// key check
	if isFound {

		// checking key only
		if len(part) == 1 {
			return true
		}

		// value check
		if val == part[1] {
			return true
		}
	}

	return false
}
