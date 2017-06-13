package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/thomersch/gosmparse"
)

// ConfigJSON - struct representing the config file format
type ConfigJSON map[string]Group

// Group - a collection of groups
type Group [][]string

// Pattern - a group if one or more AND conditions
type Pattern []string

// FeatureSet struct
type FeatureSet struct {
	NodeConditions     Group
	WayConditions      Group
	RelationConditions Group
}

// NewFeatureSetFromJSON - load a featureset from JSON
func NewFeatureSetFromJSON(path string) (*FeatureSet, error) {

	file, e := ioutil.ReadFile(path)
	if nil != e {
		return nil, e
	}

	decoder := json.NewDecoder(bytes.NewReader(file))
	var config ConfigJSON
	decoder.Decode(&config)

	var fs = &FeatureSet{}
	fs.NodeConditions = config["node"]
	fs.WayConditions = config["way"]
	fs.RelationConditions = config["relation"]

	return fs, nil
}

// MatchNode - yes/no if any target feature matches this node
func (fs *FeatureSet) MatchNode(n gosmparse.Node) bool {
	return matchAny(n.Tags, fs.NodeConditions)
}

// MatchWay - yes/no if any target feature matches this way
func (fs *FeatureSet) MatchWay(w gosmparse.Way) bool {
	return matchAny(w.Tags, fs.WayConditions)
}

// MatchRelation - yes/no if any target feature matches this relation
func (fs *FeatureSet) MatchRelation(r gosmparse.Relation) bool {
	return matchAny(r.Tags, fs.RelationConditions)
}

// matchAny - match any group (logical OR)
func matchAny(tags map[string]string, groups Group) bool {

	// no tags at all
	if len(tags) == 0 {
		return false
	}

	// trim all keys/value of extra whitespace
	tags = trimTags(tags)

	// OR groups
	for _, group := range groups {
		// AND conditions
		if matchAll(tags, group) {
			return true
		}
	}

	return false
}

// matchAll - match all conditions (logical AND)
func matchAll(tags map[string]string, pattern Pattern) bool {
	if len(pattern) == 0 {
		return false
	}
	for _, condition := range pattern {
		if !matchOne(tags, condition) {
			return false
		}
	}
	return true
}

// matchOne - match a single condition
func matchOne(tags map[string]string, condition string) bool {

	// split to array, first part will be the key and second the value (if required)
	part := strings.Split(condition, "=")
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

func trimTags(tags map[string]string) map[string]string {
	trimmed := make(map[string]string)
	for k, v := range tags {
		trimmed[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return trimmed
}
