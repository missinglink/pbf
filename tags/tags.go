package tags

// Discardable tags
// ref: http://wiki.openstreetmap.org/wiki/Discardable_tags
// ref: https://github.com/openstreetmap/iD/blob/master/data/discarded.json
func Discardable() map[string]bool {
	tags := make(map[string]bool)
	tags["created_by"] = false
	tags["converted_by"] = false
	tags["geobase:datasetName"] = false
	tags["geobase:uuid"] = false
	tags["KSJ2:ADS"] = false
	tags["KSJ2:ARE"] = false
	tags["KSJ2:AdminArea"] = false
	tags["KSJ2:COP_label"] = false
	tags["KSJ2:DFD"] = false
	tags["KSJ2:INT"] = false
	tags["KSJ2:INT_label"] = false
	tags["KSJ2:LOC"] = false
	tags["KSJ2:LPN"] = false
	tags["KSJ2:OPC"] = false
	tags["KSJ2:PubFacAdmin"] = false
	tags["KSJ2:RAC"] = false
	tags["KSJ2:RAC_label"] = false
	tags["KSJ2:RIC"] = false
	tags["KSJ2:RIN"] = false
	tags["KSJ2:WSC"] = false
	tags["KSJ2:coordinate"] = false
	tags["KSJ2:curve_id"] = false
	tags["KSJ2:curve_type"] = false
	tags["KSJ2:filename"] = false
	tags["KSJ2:lake_id"] = false
	tags["KSJ2:lat"] = false
	tags["KSJ2:long"] = false
	tags["KSJ2:river_id"] = false
	tags["odbl"] = false
	tags["odbl:note"] = false
	tags["SK53_bulk:load"] = false
	tags["sub_sea:type"] = false
	tags["tiger:source"] = false
	tags["tiger:separated"] = false
	tags["tiger:tlid"] = false
	tags["tiger:upload_uuid"] = false
	tags["yh:LINE_NAME"] = false
	tags["yh:LINE_NUM"] = false
	tags["yh:STRUCTURE"] = false
	tags["yh:TOTYUMONO"] = false
	tags["yh:TYPE"] = false
	tags["yh:WIDTH"] = false
	tags["yh:WIDTH_RANK"] = false
	return tags
}

// Uninteresting tags
// ref: http://josm.openstreetmap.de/browser/josm/trunk/src/org/openstreetmap/josm/data/osm/OsmPrimitive.java#L659
func Uninteresting() map[string]bool {
	tags := make(map[string]bool)
	tags["source"] = true
	tags["comment"] = false
	tags["note"] = false
	tags["watch"] = true
	tags["description"] = false
	tags["attribution"] = false
	return tags
}

// Highway tags relevant for finding cross streets
func Highway() map[string]bool {
	tags := make(map[string]bool)
	tags["motorway"] = false
	tags["trunk"] = false
	tags["primary"] = false
	tags["secondary"] = false
	tags["residential"] = false
	tags["service"] = false
	tags["tertiary"] = false
	tags["road"] = false
	return tags
}
