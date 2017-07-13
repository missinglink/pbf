package handler

import (
	"log"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/tags"
	"github.com/mmcloughlin/geohash"
)

// DenormalizedJSON - JSON
type DenormalizedJSON struct {
	Writer          *lib.BufferedWriter
	Conn            *leveldb.Connection
	ComputeCentroid bool
	ComputeGeohash  bool
	ExportLatLons   bool
}

// ReadNode - called once per node
func (d *DenormalizedJSON) ReadNode(item gosmparse.Node) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// node
	obj := json.NodeFromParser(item)

	// compute geohash
	if d.ComputeGeohash {
		obj.Hash = geohash.Encode(item.Lat, item.Lon)
	}

	d.Writer.Queue <- obj.Bytes()
}

// ReadWay - called once per way
func (d *DenormalizedJSON) ReadWay(item gosmparse.Way) {

	// discard selected tags
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// collect dependant node refs from store
	var refs = make([]*gosmparse.Node, 0, len(item.NodeIDs))
	for _, ref := range item.NodeIDs {
		var node, readError = d.Conn.ReadCoord(ref)
		if nil != readError {
			// skip ways which fail to denormalize
			log.Printf("skipping way %d. failed to load ref %d\n", item.ID, ref)
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

	// compute geohash
	if d.ComputeGeohash {
		obj.Hash = geohash.Encode(obj.Centroid.Lat, obj.Centroid.Lon)
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
	item.Tags = tags.Trim(item.Tags)
	DeleteTags(item.Tags, discardableTags)
	DeleteTags(item.Tags, uninterestingTags)

	// relation
	obj := json.DenormalizedRelation{
		ID:   item.ID,
		Type: "relation",
		Tags: item.Tags,
	}

	// compute polygon centroid
	if d.ComputeCentroid {

		// iterate members once to try to classify the relation
		var adminCentreID int64
		var wayIDs []int64

		for _, member := range item.Members {
			switch member.Type {
			case gosmparse.NodeType:
				// only target the 'admin_centre' node
				if member.Role == "admin_centre" {

					// store the ID of the admin centre node
					adminCentreID = member.ID
				}
			case gosmparse.WayType:
				// skip cyclic references to parent
				if member.Role != "subarea" {

					// append way ID to list of member ways
					wayIDs = append(wayIDs, member.ID)
				}
			}
		}

		// this is the simplest relation to build, we simply need to load the
		// admin centre coord and use that as the centroid
		if 0 != adminCentreID {

			var node, readError = d.Conn.ReadCoord(adminCentreID)
			if nil != readError {
				// skip relation if the point is not found in the db
				log.Printf("skipping relation %d. failed to load admin centre %d\n", item.ID, adminCentreID)
				return
			}

			// set the centroid
			obj.Centroid = json.NewLatLon(node.Lat, node.Lon)

		} else {
			// this is more complex, we need to load all the multipolygon rings
			// from the DB and assemble the geometry before calculating the centroid

			// load ring data from database
			var ways = make(map[int64]*json.DenormalizedWay)
			for _, wayID := range wayIDs {

				// load way from DB
				var way, readError = d.Conn.ReadPath(wayID)
				if nil != readError {
					// skip ways which fail to denormalize
					log.Printf("skipping relation %d. failed to load way %d\n", item.ID, wayID)
					return
				}

				// use a struct which allows us to store the refs within
				var denormalizedWay = json.DenormalizedWayFromParser(*way)

				// load way refs from DB
				for _, ref := range way.NodeIDs {
					var node, readError = d.Conn.ReadCoord(ref)
					if nil != readError {
						// skip ways which fail to denormalize
						log.Printf("skipping relation way %d. failed to load ref %d\n", item.ID, ref)
						return
					}

					// append way vertex
					denormalizedWay.LatLons = append(denormalizedWay.LatLons, json.NewLatLon(node.Lat, node.Lon))
				}

				// store way
				ways[item.ID] = denormalizedWay
			}

			log.Println("write relation", item.ID)
		}
	}

	// compute geohash
	if d.ComputeGeohash {
		obj.Hash = geohash.Encode(obj.Centroid.Lat, obj.Centroid.Lon)
	}

	d.Writer.Queue <- obj.Bytes()
}
