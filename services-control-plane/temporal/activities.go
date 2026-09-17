package temporal

import (
	"context"
	"fmt"
	"net"
	"os/exec"

	"github.com/ahmedhesham301/autoscaling-hetzner/services-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/services-control-plane/utils"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/random"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"go.temporal.io/sdk/activity"
)

func checkImageExist(ctx context.Context, params data.CreateServiceParams) (*int64, error) {
	images, err := hetzner.HClient.Image.AllWithOpts(ctx, hcloud.ImageListOpts{
		LabelSelector: services.ConvertToHetznerLabels(
			params.GetConfigMap(),
		),
	})

	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, nil
	}
	return &images[0].ID, nil

}

type buildImageParams struct {
	ServiceParams data.CreateServiceParams
	Env           string
	TemplatesPath string
	NetworkID     int64
}

func buildImage(ctx context.Context, params buildImageParams) (*int64, error) {
	args := []string{
		"build", "-machine-readable",
		"-var", fmt.Sprintf("config=%v", utils.ConvertMapToJsonString(params.ServiceParams.GetConfigMap())),
		"-var", fmt.Sprintf("env=%v", params.Env),
		"-var", fmt.Sprintf("networkID=%v", params.NetworkID),
	}

	args = append(args, params.TemplatesPath+"/"+params.ServiceParams.ServiceType+"/"+params.ServiceParams.Engine+"/main.pkr.hcl")

	logger := activity.GetLogger(ctx)
	cmd := exec.CommandContext(ctx, "packer", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("packer build command failed", "err", err, "output", string(output))
		return nil, err
	}
	logger.Info("packer build command output", "output", string(output))

	id, err := utils.GetSnapshotID(output)
	return &id, err
}

type deployServiceParams struct {
	ServiceParams      data.CreateServiceParams
	ImageID            int64
	Env                string
	AllowAllFirewallID *int64
	NetworkID          int64
}

func deployService(ctx context.Context, params deployServiceParams) error {
	ops := hcloud.ServerCreateOpts{
		Name:       random.AddRandomLetters(params.ServiceParams.Engine + "-" + params.ServiceParams.Version),
		ServerType: &hcloud.ServerType{Name: params.ServiceParams.ServerType},
		Image:      &hcloud.Image{ID: params.ImageID},
		Location:   &hcloud.Location{Name: params.ServiceParams.Location},
		PublicNet: &hcloud.ServerCreatePublicNet{
			EnableIPv4: params.ServiceParams.Network.PublicIPv4,
			EnableIPv6: params.ServiceParams.Network.PublicIPv6,
		},
		Labels: services.AppendManagedLabel(params.ServiceParams.GetConfigMapString()),
	}
	var firewalls []*hcloud.ServerCreateFirewall
	if params.ServiceParams.FirewallIDs != nil {
		for _, id := range *params.ServiceParams.FirewallIDs {
			firewalls = append(firewalls, &hcloud.ServerCreateFirewall{
				Firewall: hcloud.Firewall{
					ID: id,
				},
			})
		}
	}

	if params.Env == "dev" {
		ops.PublicNet.EnableIPv4 = true
		ops.PublicNet.EnableIPv6 = true
		firewalls = append(firewalls, &hcloud.ServerCreateFirewall{
			Firewall: hcloud.Firewall{
				ID: *params.AllowAllFirewallID,
			},
		})
	}
	ops.Firewalls = firewalls
	if params.ServiceParams.Network.PrivateNetwork {
		ops.Networks = []*hcloud.Network{
			{

				ID: params.NetworkID,
			},
		}
	}
	server, _, err := hetzner.HClient.Server.Create(ctx, ops)

	if err != nil {
		return err
	}
	return params.ServiceParams.SaveToDB(ctx, *server.Server)

}

func GetOrCreateAllowAllFirewall(ctx context.Context) (*int64, error) {
	firewall, _, err := hetzner.HClient.Firewall.GetByName(ctx, "allow_all")
	if err != nil {
		return nil, err
	}
	if firewall == nil {
		allIPv4 := net.IPNet{
			IP:   net.IPv4zero,
			Mask: net.CIDRMask(0, 32),
		}
		opts := hcloud.FirewallCreateOpts{
			Name: "allow_all",
			Rules: []hcloud.FirewallRule{
				{
					Direction: "in",
					SourceIPs: []net.IPNet{allIPv4},
					Port:      new("any"),
					Protocol:  hcloud.FirewallRuleProtocolTCP,
				},
				{
					Direction:      "out",
					DestinationIPs: []net.IPNet{allIPv4},
					Port:           new("any"),
					Protocol:       hcloud.FirewallRuleProtocolTCP,
				},
			},
		}
		firewallResult, _, err := hetzner.HClient.Firewall.Create(ctx, opts)
		if err != nil {
			return nil, err
		}
		firewall = firewallResult.Firewall
	}
	return &firewall.ID, nil
}
