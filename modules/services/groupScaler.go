package services

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/model"
	"context"
	"log/slog"
	"math/rand/v2"
	"net"
	"os"

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

type ScaleOps struct {
	GroupId int
	Group   *model.Group
}

func ScaleUp(ops ScaleOps, amount int, source string) error {
	var group model.Group
	if ops.Group == nil {
		err := group.GetById(ops.GroupId)
		if err != nil {
			return err
		}
	} else {
		group = *ops.Group
	}

	if group.DesiredSize >= group.MaxSize {
		return nil
	}

	var template model.Template
	err := template.GetById(group.TemplateId)
	if err != nil {
		return err
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

	var firewalls []*hcloud.ServerCreateFirewall
	if template.Firewalls != nil {
		for _, v := range template.Firewalls {
			firewalls = append(firewalls, &hcloud.ServerCreateFirewall{Firewall: hcloud.Firewall{ID: v}})
		}
	}

	for i := 0; i < amount; i++ {
		if group.DesiredSize >= group.MaxSize {
			return nil
		}
		whereToScale, err := whereToScaleUp("balanced", group.Locations, group.Id)
		if err != nil {
			return err
		}

		res, _, err := hetzner.HClient.Server.Create(context.Background(), hcloud.ServerCreateOpts{
			Name:       addRandomLetters(group.Name),
			ServerType: &hcloud.ServerType{Name: group.ServerType},
			Image:      &hcloud.Image{ID: template.ImageId},
			Location:   &hcloud.Location{ID: *whereToScale},
			Networks:   networks,
			UserData:   template.CloudConfig,
			SSHKeys:    SSHKeys,
			PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: *template.PublicIPv4, EnableIPv6: *template.PublicIPv6},
			Firewalls:  firewalls,
		})
		if err != nil {
			return err
		}

		// use public ip if the environment is dev
		var privateIp net.IP
		if os.Getenv("ENV") == "dev" {
			privateIp = res.Server.PublicNet.IPv4.IP
		} else {
			privateIp = res.Server.PrivateNet[0].IP
		}

		server := model.Server{
			ID:        res.Server.ID,
			GroupId:   group.Id,
			Name:      res.Server.Name,
			Type:      res.Server.ServerType.Name,
			Location:  res.Server.Location.ID,
			PrivateIp: privateIp,
		}
		err = server.Save()
		if err != nil {
			return err
		}

		if source == "alert" {

			err = group.UpdateDesiredSize(1)
			if err != nil {
				return err
			}
		}
	}

	if source == "init" {
		_, err = SetupAlert(&group)
		if err != nil {
			return err
		}
	}

	return nil
}

func ScaleOut(ops ScaleOps) error {
	var group model.Group
	if ops.Group == nil {
		err := group.GetById(ops.GroupId)
		if err != nil {
			return err
		}
	} else {
		group = *ops.Group
	}

	if group.DesiredSize <= group.MinSize {
		return nil
	}

	servers, err := model.GetAllServersInGroup(group.Id)
	if err != nil {
		return err
	}

	scaleLocationID := int64(-1)
	scaleLocationServerCount := -1
	for locationID, servers := range servers {
		if len(servers) > scaleLocationServerCount {
			scaleLocationID = locationID
			scaleLocationServerCount = len(servers)
		}
	}
	_, _, err = hetzner.HClient.Server.DeleteWithResult(context.TODO(), &hcloud.Server{ID: servers[scaleLocationID][0].ID})
	if err != nil {
		return err
	}

	err = servers[scaleLocationID][0].DeleteServer()
	if err != nil {
		return err
	}

	err = group.UpdateDesiredSize(-1)
	if err != nil {
		return err
	}
	slog.Info("group scaled out", "groupID", group.Id)
	return nil
}

// returns ID of location to start the server at
// if method is random it returns random location
// if method is something it returns the id of the location with least servers
func whereToScaleUp(method string, locations []int64, groupID int) (*int64, error) {
	if method == "random" {
		return &locations[rand.IntN(len(locations))], nil
	}
	servers, err := model.GetAllServersInGroup(groupID)
	if err != nil {
		return nil, err
	}

	// check if a location has no servers
	// if yes it return that location
	if len(locations) > len(servers) {
		for _, v := range locations {
			if _, ok := servers[v]; !ok {
				return &v, nil
			}
		}
	}

	scaleLocationID := int64(-1)
	scaleLocationServerCount := int64(999999999999999)
	for locationID, servers := range servers {
		if len(servers) < int(scaleLocationServerCount) {
			scaleLocationID = locationID
			scaleLocationServerCount = int64(len(servers))
		}
	}

	return &scaleLocationID, nil
}
