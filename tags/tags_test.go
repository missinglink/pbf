package tags

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscardable(t *testing.T) {
	var tagMap = Discardable()
	assert.Equal(t, "map[string]bool", reflect.TypeOf(tagMap).String())
	assert.True(t, len(tagMap) > 0)
}

func TestUninteresting(t *testing.T) {
	var tagMap = Uninteresting()
	assert.Equal(t, "map[string]bool", reflect.TypeOf(tagMap).String())
	assert.True(t, len(tagMap) > 0)
}

func TestHighway(t *testing.T) {
	var tagMap = Highway()
	assert.Equal(t, "map[string]bool", reflect.TypeOf(tagMap).String())
	assert.True(t, len(tagMap) > 0)
}
