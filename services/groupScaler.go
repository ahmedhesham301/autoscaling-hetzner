package services

import (
	"autoscaling-hetzner/hetzner"
	"autoscaling-hetzner/model"
	"context"
	"errors"
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
	if len(group.Locations) == 0 {
		return errors.New("group has no locations configured")
	}

	var networks []*hcloud.Network
	if template.Networks != nil {
		for _, v := range template.Networks {
			networks = append(networks, &hcloud.Network{ID: v})
		}
	}

	var SSHKeys []*hcloud.SSHKey
	for _, v := range template.SSHKeys {
		SSHKeys = append(SSHKeys, &hcloud.SSHKey{ID: v})
	}
	for i := 0; i < amount; i++ {
		res, _, err := hetzner.HClient.Server.Create(context.Background(), hcloud.ServerCreateOpts{
			Name:       addRandomLetters(group.Name),
			ServerType: &hcloud.ServerType{Name: group.ServerType},
			Image:      &hcloud.Image{ID: template.ImageId},
			Location:   &hcloud.Location{ID: group.Locations[i%len(group.Locations)]},
			Networks:   networks,
			UserData:   template.CloudConfig,
			SSHKeys:    SSHKeys,
			PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: *template.PublicIPv4, EnableIPv6: *template.PublicIPv6},
		})
		if err != nil {
			return err
		}

		server := model.Server{
			GroupId:   group.Id,
			Name:      res.Server.Name,
			Type:      res.Server.ServerType.Name,
			Location:  res.Server.Location.ID,
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
