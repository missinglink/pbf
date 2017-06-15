package handler

import (
	"log"
	"sync"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
)

// DenormlizedJSON - JSON
type DenormlizedJSON struct {
	Mutex *sync.Mutex
	Conn  *leveldb.Connection
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
		var node, readError = d.Conn.ReadNode(ref)
		if nil != readError {
			log.Printf("failed to load noderef: %d\n", ref)

			// skip ways which fail to denormalize
			return
		} else {
			refs = append(refs, node)
		}
	}

	// compute line/street centroid
	var lon, lat = lib.WayCentroid(refs)

	// convert refs to latlons
	var latlons []*json.LatLon
	for _, node := range refs {
		latlons = append(latlons, &json.LatLon{
			Lat: node.Lat,
			Lon: node.Lon,
		})
	}

	// way
	json := json.DernomalizedWay{
		ID:       item.ID,
		Type:     "way",
		Tags:     item.Tags,
		Centroid: json.NewLatLon(lat, lon),
		LatLons:  latlons,
	}

	d.Mutex.Lock()
	json.Print()
	d.Mutex.Unlock()
}

// ReadRelation - called once per relation
func (d *DenormlizedJSON) ReadRelation(item gosmparse.Relation) {
	/* currently unsupported */
}
