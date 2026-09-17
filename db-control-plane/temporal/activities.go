package temporal

import (
	"context"
	"fmt"
	"net"
	"os/exec"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/utils"
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
	serviceParams data.CreateServiceParams
	env           string
	templatesPath string
	networkID     int64
}

func buildImage(ctx context.Context, params buildImageParams) (*int64, error) {
	args := []string{
		"build", "-machine-readable",
		"-var", fmt.Sprintf("config=%v", utils.ConvertMapToJsonString(params.serviceParams.GetConfigMap())),
		"-var", fmt.Sprintf("env=%v", params.env),
		"-var", fmt.Sprintf("networkID=%v", params.networkID),
	}

	args = append(args, params.templatesPath+"/"+params.serviceParams.ServiceType+"/"+params.serviceParams.Engine+"/main.pkr.hcl")

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
	serviceParams      data.CreateServiceParams
	imageID            int64
	env                string
	allowAllFirewallID *int64
	networkID          int64
}

func deployService(ctx context.Context, params deployServiceParams) error {
	ops := hcloud.ServerCreateOpts{
		Name:       random.AddRandomLetters(params.serviceParams.Engine + "-" + params.serviceParams.Version),
		ServerType: &hcloud.ServerType{Name: params.serviceParams.ServerType},
		Image:      &hcloud.Image{ID: params.imageID},
		Location:   &hcloud.Location{Name: params.serviceParams.Location},
		PublicNet: &hcloud.ServerCreatePublicNet{
			EnableIPv4: params.serviceParams.Network.PublicIPv4,
			EnableIPv6: params.serviceParams.Network.PublicIPv6,
		},
		Labels: services.AppendManagedLabel(params.serviceParams.GetConfigMapString()),
	}
	var firewalls []*hcloud.ServerCreateFirewall
	for _, id := range *params.serviceParams.FirewallIDs {
		firewalls = append(firewalls, &hcloud.ServerCreateFirewall{
			Firewall: hcloud.Firewall{
				ID: id,
			},
		})
	}

	if params.env == "dev" {
		ops.PublicNet.EnableIPv4 = true
		ops.PublicNet.EnableIPv6 = true
		firewalls = append(firewalls, &hcloud.ServerCreateFirewall{
			Firewall: hcloud.Firewall{
				ID: *params.allowAllFirewallID,
			},
		})
	}
	ops.Firewalls = firewalls
	if params.serviceParams.Network.PrivateNetwork {
		ops.Networks = []*hcloud.Network{
			{

				ID: params.networkID,
			},
		}
	}
	server, _, err := hetzner.HClient.Server.Create(ctx, ops)

	if err != nil {
		return err
	}
	return params.serviceParams.SaveToDB(ctx, *server.Server)

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
