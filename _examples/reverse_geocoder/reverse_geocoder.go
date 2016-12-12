package main

import (
	"fmt"
	"log"

	yolp "github.com/ngs/go-yolp"
)

func main() {
	client, err := yolp.NewFromEnvionment()
	if err != nil {
		log.Fatal(err)
	}
	req := client.ReverseGeocoder(yolp.GeocoderParams{
		Latitude:  35.62172852580437,
		Longitude: 139.6999476850032,
		Datum:     yolp.WGS,
	})
	res, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Feature[0].Property.Address)
}
