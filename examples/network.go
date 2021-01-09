package main

import (
	"github.com/goccy/go-json"
	"github.com/sujit-baniya/utils/xopen"
	"sync"
)

type Network struct {
	Plmn        string `json:"plmn"`
	NibbledPlmn string `json:"nibbledPlmn"`
	Mcc         string `json:"mcc"`
	Mnc         string `json:"mnc"`
	Region      string `json:"region"`
	Type        string `json:"type"`
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	Lat         string `json:"lat"`
	Long        string `json:"long"`
	Brand       string `json:"brand"`
	Operator    string `json:"operator"`
	Status      string `json:"status"`
	Bands       string `json:"bands"`
	Notes       string `json:"notes"`
}

type NetMap []Network

var CountryNetwork = map[string][]Network{}

func LoadNetworks(path string) {
	var err error
	var pool = sync.Pool{
		New: func() interface{} { return new(NetMap) },
	}
	r, err := xopen.Ropen(path)
	if err != nil {
		panic(err)
	}
	items := pool.Get().(*NetMap)
	decoder := json.NewDecoder(r)
	_ = decoder.Decode(&items)
	for _, mp := range *items {
		CountryNetwork[mp.CountryCode] = append(CountryNetwork[mp.CountryCode], mp)
	}
	pool.Put(items)
}
