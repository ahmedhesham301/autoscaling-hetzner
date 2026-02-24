package vars

import (
	"autoscaling-hetzner/hetzner"
	"context"
	"log"
)

var Zones map[string]map[string]int64

func GetZones() {
	resp, err := hetzner.HClient.Location.All(context.Background())
	if err != nil {
		log.Panic(err.Error())
	}

	Zones = make(map[string]map[string]int64)
	for _, location := range resp {
		if Zones[string(location.NetworkZone)] == nil {
			Zones[string(location.NetworkZone)] = make(map[string]int64)
		}
		Zones[string(location.NetworkZone)][location.City] = location.ID
	}
}
