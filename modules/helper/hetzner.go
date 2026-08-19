package helper

import (
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func ToServerCreateOpts(s schema.ServerCreateRequest) (error, *hcloud.ServerCreateOpts) {
	opts := hcloud.ServerCreateOpts{
		Name:             s.Name,
		UserData:         s.UserData,
		StartAfterCreate: s.StartAfterCreate,
		Automount:        s.Automount,
	}
	// Labels
	if s.Labels != nil {
		opts.Labels = services.AppendManagedLabel(*s.Labels)
	}

	// Server type
	opts.ServerType = &hcloud.ServerType{
		ID:   s.ServerType.ID,
		Name: s.ServerType.Name,
	}

	// Image
	opts.Image = &hcloud.Image{
		ID:   s.Image.ID,
		Name: s.Image.Name,
	}

	// SSH keys
	for _, id := range s.SSHKeys {
		opts.SSHKeys = append(opts.SSHKeys, &hcloud.SSHKey{
			ID: id,
		})
	}

	// Location
	if s.Location != "" {
		if id, err := strconv.ParseInt(s.Location, 10, 64); err == nil {
			opts.Location = &hcloud.Location{
				ID: id,
			}
		} else {
			opts.Location = &hcloud.Location{
				Name: s.Location,
			}
		}
	}
	println(s.Location)
	// Volumes
	for _, id := range s.Volumes {
		opts.Volumes = append(opts.Volumes, &hcloud.Volume{
			ID: id,
		})
	}

	// Networks
	for _, id := range s.Networks {
		opts.Networks = append(opts.Networks, &hcloud.Network{
			ID: id,
		})
	}

	// Firewalls
	for _, firewall := range s.Firewalls {
		opts.Firewalls = append(opts.Firewalls, &hcloud.ServerCreateFirewall{
			Firewall: hcloud.Firewall{
				ID: firewall.Firewall,
			},
		})
	}

	// Placement group
	if s.PlacementGroup != 0 {
		opts.PlacementGroup = &hcloud.PlacementGroup{
			ID: s.PlacementGroup,
		}
	}

	// Public network
	if s.PublicNet != nil {
		opts.PublicNet = &hcloud.ServerCreatePublicNet{
			EnableIPv4: s.PublicNet.EnableIPv4,
			EnableIPv6: s.PublicNet.EnableIPv6,
		}
	}

	if err := opts.Validate(); err != nil {
		return err, nil
	}

	return nil, &opts
}
