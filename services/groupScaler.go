package services

import (
	"autoscaling-hetzner/hetzner"
	"autoscaling-hetzner/model"
	"autoscaling-hetzner/vars"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// TODO use strings.Builder
func addRandomLetters(name string) string {
	name += "-"
	for i := 0; i < 5; i++ {
		name += string(chars[rand.IntN(len(chars))])
	}
	return name
}

type ScaleUpOps struct {
	GroupId int
	Group   *model.Group
}

func ScaleUp(ops ScaleUpOps, amount int) error {
	var group model.Group
	if ops.Group == nil {
		err := group.GetById(ops.GroupId)
		if err != nil {
			return err
		}
	} else {
		group = *ops.Group
	}

	var template model.Template
	err := template.GetById(group.TemplateId)
	if err != nil {
		return err
	}
	serverTypes := group.ServerTypes

	for i := 0; i < amount; i++ {
		locationId := vars.Zones[group.Zone][group.Locations[i % len(group.Locations)]]
		if locationId == 0{
			return errors.New("There is a problem with the zone/location")	
		}
		res, _, err := hetzner.HClient.Server.Create(context.Background(), hcloud.ServerCreateOpts{
			Name:       addRandomLetters(group.Name),
			ServerType: &hcloud.ServerType{Name: serverTypes[0]},
			Image:      &hcloud.Image{Name: fmt.Sprintf("%s-%s", template.OSFlavor, template.OSVersion)},
			Location:   &hcloud.Location{ID: locationId},
			Networks:   []*hcloud.Network{&hcloud.Network{ID: 11952339}},
			UserData:   template.CloudConfig,
		})
		if err != nil {
			return err
		}

		server := model.Server{
			GroupId:   group.Id,
			Name:      res.Server.Name,
			Type:      res.Server.ServerType.Name,
			Location:  res.Server.Location.City,
			PrivateIp: res.Server.PrivateNet[0].IP,
		}
		err = server.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func ScaleOut() {

}
