package lib

import (
	"testing"

	"github.com/paulmach/go.geo"
	"github.com/stretchr/testify/assert"
)

// http://www.openstreetmap.org/way/46340228
func TestGetLineCentroid(t *testing.T) {

	var poly = geo.NewPointSet()
	poly.Push(geo.NewPoint(-74.001559, 40.719743))
	poly.Push(geo.NewPoint(-73.999914, 40.721679))
	poly.Push(geo.NewPoint(-73.997783, 40.724195))
	poly.Push(geo.NewPoint(-73.997318, 40.724745))
	poly.Push(geo.NewPoint(-73.996797, 40.725375))
	poly.Push(geo.NewPoint(-73.995203, 40.727239))
	poly.Push(geo.NewPoint(-73.993927, 40.728737))
	poly.Push(geo.NewPoint(-73.992407, 40.730535))
	poly.Push(geo.NewPoint(-73.991545, 40.731566))
	poly.Push(geo.NewPoint(-73.991417, 40.731843))
	poly.Push(geo.NewPoint(-73.990745, 40.734738))
	poly.Push(geo.NewPoint(-73.990199, 40.737495))
	poly.Push(geo.NewPoint(-73.989630, 40.739735))
	poly.Push(geo.NewPoint(-73.989370, 40.741459))
	poly.Push(geo.NewPoint(-73.989219, 40.742233))
	poly.Push(geo.NewPoint(-73.989119, 40.743025))
	poly.Push(geo.NewPoint(-73.988699, 40.745262))
	poly.Push(geo.NewPoint(-73.987904, 40.749446))
	poly.Push(geo.NewPoint(-73.987417, 40.752149))
	poly.Push(geo.NewPoint(-73.986938, 40.754016))
	poly.Push(geo.NewPoint(-73.986833, 40.754345))
	poly.Push(geo.NewPoint(-73.986321, 40.755897))
	poly.Push(geo.NewPoint(-73.986117, 40.756513))
	poly.Push(geo.NewPoint(-73.985720, 40.757348))
	poly.Push(geo.NewPoint(-73.985433, 40.757980))
	poly.Push(geo.NewPoint(-73.983607, 40.760503))
	poly.Push(geo.NewPoint(-73.979957, 40.765504))
	poly.Push(geo.NewPoint(-73.979264, 40.766480))

	var centroid = GetLineCentroid(poly)
	assert.Equal(t, 40.74239780132512, centroid.Lat())
	assert.Equal(t, -73.98919819175188, centroid.Lng())
}

// http://www.openstreetmap.org/way/264768896
func TestGetPolygonCentroid(t *testing.T) {

	var poly = geo.NewPointSet()
	poly.Push(geo.NewPoint(-73.989605, 40.740760))
	poly.Push(geo.NewPoint(-73.989615, 40.740762))
	poly.Push(geo.NewPoint(-73.989619, 40.740763))
	poly.Push(geo.NewPoint(-73.989855, 40.740864))
	poly.Push(geo.NewPoint(-73.989859, 40.740867))
	poly.Push(geo.NewPoint(-73.989866, 40.740874))
	poly.Push(geo.NewPoint(-73.989870, 40.740882))
	poly.Push(geo.NewPoint(-73.989872, 40.740891))
	poly.Push(geo.NewPoint(-73.989870, 40.740899))
	poly.Push(geo.NewPoint(-73.989865, 40.740907))
	poly.Push(geo.NewPoint(-73.989584, 40.741288))
	poly.Push(geo.NewPoint(-73.989575, 40.741294))
	poly.Push(geo.NewPoint(-73.989564, 40.741298))
	poly.Push(geo.NewPoint(-73.989559, 40.741300))
	poly.Push(geo.NewPoint(-73.989547, 40.741300))
	poly.Push(geo.NewPoint(-73.989535, 40.741299))
	poly.Push(geo.NewPoint(-73.989529, 40.741297))
	poly.Push(geo.NewPoint(-73.989519, 40.741293))
	poly.Push(geo.NewPoint(-73.989514, 40.741290))
	poly.Push(geo.NewPoint(-73.989507, 40.741283))
	poly.Push(geo.NewPoint(-73.989501, 40.741265))
	poly.Push(geo.NewPoint(-73.989570, 40.740776))
	poly.Push(geo.NewPoint(-73.989575, 40.740770))
	poly.Push(geo.NewPoint(-73.989581, 40.740765))
	poly.Push(geo.NewPoint(-73.989590, 40.740761))
	poly.Push(geo.NewPoint(-73.989595, 40.740760))
	poly.Push(geo.NewPoint(-73.989605, 40.740760))

	var centroid = GetPolygonCentroid(poly)
	assert.Equal(t, 40.74100992600508, centroid.Lat())
	assert.Equal(t, -73.98964244467275, centroid.Lng())
}
