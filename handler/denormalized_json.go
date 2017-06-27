package handler

import (
	"log"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
)

// DenormalizedJSON - JSON
type DenormalizedJSON struct {
	Writer          *lib.BufferedWriter
	Conn            *leveldb.Connection
	ComputeCentroid bool
	ExportLatLons   bool
}

// ReadNode - called once per node
func (d *DenormalizedJSON) ReadNode(item gosmparse.Node) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// node
	obj := json.NodeFromParser(item)
	d.Writer.Queue <- obj.Bytes()
}

// ReadWay - called once per way
func (d *DenormalizedJSON) ReadWay(item gosmparse.Way) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// collect dependant node refs from store
	var refs []*gosmparse.Node
	for _, ref := range item.NodeIDs {
		var node, readError = d.Conn.ReadNode(ref)
		if nil != readError {
			// skip ways which fail to denormalize
			log.Printf("failed to load noderef: %d\n", ref)
			return
		}
		refs = append(refs, node)
	}

	// way
	obj := json.DenormalizedWay{
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

	// write
	d.Writer.Queue <- obj.Bytes()
}

// ReadRelation - called once per relation
func (d *DenormalizedJSON) ReadRelation(item gosmparse.Relation) {

	// discard selected tags
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// relation
	obj := json.RelationFromParser(item)
	d.Writer.Queue <- obj.Bytes()
}
