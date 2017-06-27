package parser

import (
	"fmt"
	"sync"

	"github.com/missinglink/pbf/handler"

	"github.com/missinglink/gosmparse"
)

// const MAX_CACHE_ENTRIES = 200

// RandomAccessParser - struct to handle random access lookups to a pbf
type RandomAccessParser struct {
	Parser
	Index *gosmparse.BlobIndex
	Cache *handler.ReadAll
}

// NewRandomAccessParser -
func NewRandomAccessParser(path string, idxPath string) *RandomAccessParser {

	// load index
	index := &gosmparse.BlobIndex{}
	index.ReadFromFile(idxPath)

	var p = &RandomAccessParser{
		Index: index,
		Cache: &handler.ReadAll{
			Mutex:     &sync.Mutex{},
			Nodes:     make(map[int64]gosmparse.Node),
			Ways:      make(map[int64]gosmparse.Way),
			Relations: make(map[int64]gosmparse.Relation),
		},
	}

	p.open(path)
	return p
}

// GetNode - fetch a single record from the file
func (p *RandomAccessParser) GetNode(osmID int64) (gosmparse.Node, error) {

	// check if we have this element in the cache
	if found, ok := p.Cache.Nodes[osmID]; ok {
		return found, nil
	}

	p.loadBlob("node", osmID)

	// check if we have this element in the cache
	if found, ok := p.Cache.Nodes[osmID]; ok {
		return found, nil
	}

	return gosmparse.Node{}, fmt.Errorf("node not found: %d", osmID)
}

// GetWay - fetch a single record from the file
func (p *RandomAccessParser) GetWay(osmID int64) (gosmparse.Way, error) {

	// check if we have this element in the cache
	if found, ok := p.Cache.Ways[osmID]; ok {
		return found, nil
	}

	p.loadBlob("way", osmID)

	// check if we have this element in the cache
	if found, ok := p.Cache.Ways[osmID]; ok {
		return found, nil
	}

	return gosmparse.Way{}, fmt.Errorf("way not found: %d", osmID)
}

// GetRelation - fetch a single record from the file
func (p *RandomAccessParser) GetRelation(osmID int64) (gosmparse.Relation, error) {

	// check if we have this element in the cache
	if found, ok := p.Cache.Relations[osmID]; ok {
		return found, nil
	}

	p.loadBlob("relation", osmID)

	// check if we have this element in the cache
	if found, ok := p.Cache.Relations[osmID]; ok {
		return found, nil
	}

	return gosmparse.Relation{}, fmt.Errorf("relation not found: %d", osmID)
}

// loadBlob - fetch blob and cache returned elements
func (p *RandomAccessParser) loadBlob(osmType string, osmID int64) error {

	// find the location of this element in file
	offsets, err := p.Index.BlobOffsets(osmType, osmID)
	if nil != err {
		fmt.Printf("target element: %s %d not found in file\n", osmType, osmID)
		return err
	}

	for _, offset := range offsets {

		// Parse will block until it is done or an error occurs.
		p.ParseBlob(p.Cache, offset)

	}

	return nil
}
