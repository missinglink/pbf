package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/missinglink/gosmparse"
	"github.com/missinglink/pbf/json"
	"github.com/missinglink/pbf/leveldb"
	"github.com/missinglink/pbf/lib"
	"github.com/missinglink/pbf/spatialite"
	"github.com/missinglink/pbf/tags"
	"github.com/mmcloughlin/geohash"
)

// DenormalizedJSON - JSON
type DenormalizedJSON struct {
	Writer          *lib.BufferedWriter
	Conn            *leveldb.Connection
	Spatialite      *spatialite.Connection
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
		var nodeCentroidID int64
		var wayIDs []int64

		for _, member := range item.Members {
			switch member.Type {
			case gosmparse.NodeType:
				// only target the 'label' or 'admin_centre' nodes
				if member.Role == "label" || member.Role == "admin_centre" {

					// store the ID of the node which contains the centroid info
					nodeCentroidID = member.ID
				}
			case gosmparse.WayType:
				// skip cyclic references to parent (subarea) and other junk roles
				if member.Role == "outer" || member.Role == "inner" || member.Role == "" {

					// append way ID to list of member ways
					wayIDs = append(wayIDs, member.ID)
				}
			}
		}

		// this is the simplest relation to build, we simply need to load the
		// 'label' or 'admin_centre' node its lat/lon as the relation centroid
		if 0 != nodeCentroidID {

			var node, readError = d.Conn.ReadCoord(nodeCentroidID)
			if nil != readError {
				// skip relation if the point is not found in the db
				log.Printf("skipping relation %d. failed to load admin centre %d\n", item.ID, nodeCentroidID)
				return
			}

			// set the centroid
			obj.Centroid = json.NewLatLon(node.Lat, node.Lon)

		} else {
			// this is more complex, we need to load all the multipolygon linestrings
			// from the DB and assemble the geometry before calculating the centroid

			// generate WKT strings as input for 'GeomFromText'
			var lineStrings []string
			for _, wayID := range wayIDs {

				// load way from DB
				var way, readError = d.Conn.ReadPath(wayID)
				if nil != readError {
					// skip ways which fail to denormalize
					log.Printf("skipping relation %d. failed to load way %d\n", item.ID, wayID)
					return
				}

				// load vertices from DB
				var vertices []string
				for _, ref := range way.NodeIDs {
					var node, readError = d.Conn.ReadCoord(ref)
					if nil != readError {
						// skip ways which fail to denormalize
						log.Printf("skipping relation way %d. failed to load ref %d\n", item.ID, ref)
						return
					}

					vertices = append(vertices, fmt.Sprintf("%f %f", node.Lon, node.Lat))
				}

				lineStrings = append(lineStrings, fmt.Sprintf("(%s)", strings.Join(vertices, ",")))
			}

			// build SQL query
			var query = `SELECT COALESCE( AsText( PointOnSurface( BuildArea( GeomFromText('MULTILINESTRING(`
			query += strings.Join(lineStrings, ",")
			query += `)')))),'');`

			// query database for result
			var res string
			var err = d.Spatialite.DB.QueryRow(query).Scan(&res)
			if err != nil {
				log.Printf("spatialite: failed to assemble relation: %d", item.ID)
				return
			}

			// extract lat/lon values from WKT
			var lon, lat float64
			n, _ := fmt.Sscanf(res, "POINT(%f %f)", &lon, &lat)

			// ensure we got 2 floats
			if 2 != n {
				log.Printf("spatialite: failed to compute centroid for relation: %d", item.ID)
				return
			}

			// set the centroid
			obj.Centroid = json.NewLatLon(lat, lon)
		}
	}

	// compute geohash
	if d.ComputeGeohash {
		obj.Hash = geohash.Encode(obj.Centroid.Lat, obj.Centroid.Lon)
	}

	d.Writer.Queue <- obj.Bytes()
}
