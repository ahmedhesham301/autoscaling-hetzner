package services

import (
	"autoscaling-hetzner/hetzner"
	"autoscaling-hetzner/model"
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func ScaleUp(groupId int) error {
	var group model.Group
	print("before calling the api")
	err := group.GetById(groupId)
	if err != nil {
		return err
	}

	var template model.Template
	err = template.GetById(group.TemplateId)
	if err != nil {
		return err
	}
	res, _, err := hetzner.HClient.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       "testnameee",
		ServerType: &hcloud.ServerType{Name: "cx23"},
		Image:      &hcloud.Image{Name: fmt.Sprintf("%s-%s", template.OSFlavor, template.OSVersion)},
		Location:   &hcloud.Location{City: group.Locations[0]},
		Networks:   []*hcloud.Network{&hcloud.Network{ID: 11952339}},
	})
	if err != nil {
		return err
	}
	server := model.Server{
		GroupId:   group.Id,
		Type:      res.Server.ServerType.Name,
		Location:  res.Server.Location.Country,
		PrivateIp: res.Server.PrivateNet[0].IP,
	}
	err = server.Save()
	if err != nil {
		return err
	}
	return nil
}

func ScaleOut() {

}
