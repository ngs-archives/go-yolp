package yolp

import "testing"

func TestGeocoderParamsQuery(t *testing.T) {
	for _, test := range []Test{
		{
			"lat=35.62172852580437&lon=139.6999476850032",
			GeocoderParams{
				Latitude:  35.62172852580437,
				Longitude: 139.6999476850032,
			}.Query().Encode(),
		},
		{
			"datum=wgs&lat=37.7451301&lon=-122.5680713",
			GeocoderParams{
				Latitude:  37.7451301,
				Longitude: -122.5680713,
				Datum:     WGS,
			}.Query().Encode(),
		},
		{
			"datum=tky&lat=35.618486&lon=139.703173528",
			GeocoderParams{
				Latitude:  35.618486000,
				Longitude: 139.703173528,
				Datum:     Tokyo,
			}.Query().Encode(),
		},
	} {
		test.Compare(t)
	}
}
