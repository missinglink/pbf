package json

// full OSM precision of 7 decimal places
var precision = 10000000.0

// return a float64 truncated to the desired precision
func truncate(val float64) float64 {
	return float64(int(val*precision)) / precision
}
