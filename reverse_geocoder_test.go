package yolp

import (
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestReverseGeocoderErrorResponse(t *testing.T) {
	client, _ := New("APP", "SEC")
	op := client.ReverseGeocoder(GeocoderParams{
		Longitude: 35.62172852580437,
		Latitude:  139.6999476850032,
		Datum:     "wsg",
	})
	gock.New("http://reverse.search.olp.yahooapis.jp/OpenLocalPlatform/V1/reverseGeoCoder?datum=wsg&lat=139.6999476850032&lon=35.62172852580437&output=xml").
		Reply(200).
		File("_fixtures/error.xml")
	res, err := op.Do()
	if err == nil {
		t.Errorf("Expected not nil but got nil res: %v", res)
	} else {
		Test{"Your Request was Forbidden", err.Error()}.Compare(t)
	}
}

func TestReverseGeocoderResponse(t *testing.T) {
	client, _ := New("APP", "SEC")
	op := client.ReverseGeocoder(GeocoderParams{
		Longitude: 35.62172852580437,
		Latitude:  139.6999476850032,
		Datum:     "wsg",
	})
	gock.New("http://reverse.search.olp.yahooapis.jp/OpenLocalPlatform/V1/reverseGeoCoder?datum=wsg&lat=139.6999476850032&lon=35.62172852580437&output=xml").
		Reply(200).
		File("_fixtures/reverse_geocoder.xml")
	res, err := op.Do()
	if err != nil {
		t.Errorf("Expected nil but got %v", err.Error())
	}
	for _, test := range []Test{
		{1, res.ResultInfo.Count},
		{1, res.ResultInfo.Total},
		{1, res.ResultInfo.Start},
		{float64(0.0052351951599121), res.ResultInfo.Latency},
		{200, res.ResultInfo.Status},
		{"指定の地点の住所情報を取得する機能を提供します。", res.ResultInfo.Description},
		{"Copyright (C) 2016 Yahoo Japan Corporation. All Rights Reserved.", res.ResultInfo.Copyright},
		{"", res.ResultInfo.CompressType},
		{1, len(res.Feature)},
		{"", res.Feature[0].Description},
		{139.69994768500328, res.Feature[0].Geometry.Coordinates.Latitude},
		{35.62172852580437, res.Feature[0].Geometry.Coordinates.Longitude},
		{"", res.Feature[0].ID},
		{"", res.Feature[0].Name},
		{"東京都目黒区目黒本町５丁目１６", res.Feature[0].Property.Address},
		{5, len(res.Feature[0].Property.AddressElement)},
		{"13", res.Feature[0].Property.AddressElement[0].Code},
		{"とうきょうと", res.Feature[0].Property.AddressElement[0].Kana},
		{"prefecture", res.Feature[0].Property.AddressElement[0].Level},
		{"東京都", res.Feature[0].Property.AddressElement[0].Name},
		{"13110", res.Feature[0].Property.AddressElement[1].Code},
		{"めぐろく", res.Feature[0].Property.AddressElement[1].Kana},
		{"city", res.Feature[0].Property.AddressElement[1].Level},
		{"目黒区", res.Feature[0].Property.AddressElement[1].Name},
		{"", res.Feature[0].Property.AddressElement[2].Code},
		{"めぐろほんちょう", res.Feature[0].Property.AddressElement[2].Kana},
		{"oaza", res.Feature[0].Property.AddressElement[2].Level},
		{"目黒本町", res.Feature[0].Property.AddressElement[2].Name},
		{"", res.Feature[0].Property.AddressElement[3].Code},
		{"５ちょうめ", res.Feature[0].Property.AddressElement[3].Kana},
		{"aza", res.Feature[0].Property.AddressElement[3].Level},
		{"５丁目", res.Feature[0].Property.AddressElement[3].Name},
		{"", res.Feature[0].Property.AddressElement[4].Code},
		{"１６", res.Feature[0].Property.AddressElement[4].Kana},
		{"detail1", res.Feature[0].Property.AddressElement[4].Level},
		{"１６", res.Feature[0].Property.AddressElement[4].Name},
		{1, len(res.Feature[0].Property.Building)},
		{"test area", res.Feature[0].Property.Building[0].Area},
		{"3", res.Feature[0].Property.Building[0].Floor},
		{"B@n0WrVSt_O", res.Feature[0].Property.Building[0].ID},
		{"test name", res.Feature[0].Property.Building[0].Name},
		{1, len(res.Feature[0].Property.Road)},
		{"Test Name", res.Feature[0].Property.Road[0].Name},
		{"Test Kana", res.Feature[0].Property.Road[0].Kana},
		{"Test PopularName", res.Feature[0].Property.Road[0].PopularName},
		{"Test PopularKana", res.Feature[0].Property.Road[0].PopularKana},
	} {
		test.Compare(t)
	}
}
