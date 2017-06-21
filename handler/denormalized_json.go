package handler

import (
	"log"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/parser"
)

// DenormlizedJSON - JSON
type DenormlizedJSON struct {
	Mutex           *sync.Mutex
	Store           *parser.CachedRandomAccessParser
	ComputeCentroid bool
	ExportLatLons   bool
}

// ReadNode - called once per node
func (d *DenormlizedJSON) ReadNode(item gosmparse.Node) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// node
	json := json.Node{
		ID:   item.ID,
		Type: "node",
		Lat:  item.Lat,
		Lon:  item.Lon,
		Tags: item.Tags,
	}

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}

// ReadWay - called once per way
func (d *DenormlizedJSON) ReadWay(item gosmparse.Way) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// collect dependant node refs from store
	var refs []*gosmparse.Node
	for _, ref := range item.NodeIDs {
		var node, readError = d.Store.ReadNode(ref)
		if nil != readError {
			log.Printf("failed to load noderef: %d\n", ref)

			// skip ways which fail to denormalize
			return
		}
		refs = append(refs, node)
	}

	// way
	obj := json.DernomalizedWay{
		ID:   item.ID,
		Type: "way",
		Tags: item.Tags,
	}

	// compute line/street centroid
	if d.ComputeCentroid {
		var lon, lat = lib.WayCentroid(refs)
		obj.Centroid = json.NewLatLon(lat, lon)
	}

	// convert refs to latlons
	if d.ExportLatLons {
		for _, node := range refs {
			obj.LatLons = append(obj.LatLons, &json.LatLon{
				Lat: node.Lat,
				Lon: node.Lon,
			})
		}
	}

	d.Mutex.Lock()
	obj.Print()
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *DenormlizedJSON) ReadRelation(item gosmparse.Relation) {
	/* currently unsupported */
}
