package lib

import (
	"github.com/missinglink/gosmparse"
	"github.com/paulmach/go.geo"
)

// WayCentroid - compute the centroid of a way
func WayCentroid(refs []*gosmparse.Node) (float64, float64) {

	// convert lat/lon map to geo.PointSet
	points := geo.NewPointSet()
	for _, node := range refs {
		points.Push(geo.NewPoint(float64(node.Lon), float64(node.Lat)))
	}

	// determine if the way is a closed centroid or a linestring
	// by comparing first and last coordinates.
	isClosed := false
	if points.Length() > 2 {
		isClosed = points.First().Equals(points.Last())
	}

	// compute the centroid using one of two different algorithms
	var compute *geo.Point
	if isClosed {
		compute = GetPolygonCentroid(points)
	} else {
		compute = GetLineCentroid(points)
	}

	// return centroid
	var lat = compute.Lat()
	var lon = compute.Lng()

	return lon, lat
}

// GetPolygonCentroid - compute the centroid of a polygon set
// using a spherical co-ordinate system
func GetPolygonCentroid(ps *geo.PointSet) *geo.Point {
	// GeoCentroid function added in https://github.com/paulmach/go.geo/pull/24
	return ps.GeoCentroid()
}

// GetLineCentroid - compute the centroid of a line string
func GetLineCentroid(ps *geo.PointSet) *geo.Point {

	path := geo.NewPath()
	path.PointSet = *ps

	halfDistance := path.Distance() / 2
	travelled := 0.0

	for i := 0; i < len(path.PointSet)-1; i++ {

		segment := geo.NewLine(&path.PointSet[i], &path.PointSet[i+1])
		distance := segment.Distance()

		// middle line segment
		if (travelled + distance) > halfDistance {
			var remainder = halfDistance - travelled
			return segment.Interpolate(remainder / distance)
		}

		travelled += distance
	}

	return ps.GeoCentroid()
}
