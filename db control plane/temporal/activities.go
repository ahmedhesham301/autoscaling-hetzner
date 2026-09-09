package temporal

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/model"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/utils"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/random"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"go.temporal.io/sdk/activity"
)

func checkImageExist(ctx context.Context, params data.CreateDBParams) (*int64, error) {
	images, err := hetzner.HClient.Image.AllWithOpts(ctx, hcloud.ImageListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: services.ConvertToHetznerLabels(
				params.GetConfigMap(),
			),
		},
	})

	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, nil
	}
	return &images[0].ID, nil

}

func buildImage(ctx context.Context, params data.CreateDBParams, env string, templatesPath string) (*int64, error) {
	logger := activity.GetLogger(ctx)

	cmd := exec.CommandContext(ctx, "packer", "build", "-machine-readable",
		"-var", fmt.Sprintf("config=%v", utils.ConvertMapToJsonString(params.GetConfigMap())),
		"-var", fmt.Sprintf("ENV=%v", env),
		templatesPath+"/"+params.AppName+"/main.pkr.hcl")

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("packer build command failed", "err", err, "output", string(output))
		return nil, err
	}
	logger.Info("packer build command output", "output", string(output))

	id, err := utils.GetSnapshotID(output)
	return &id, err
}

func deployDB(ctx context.Context, params data.CreateDBParams, imageID int64, DB_ID int) error {
	env := os.Getenv("ENV")

	ops := hcloud.ServerCreateOpts{
		Name:       random.AddRandomLetters(params.AppName + "-" + params.AppVersion),
		ServerType: &hcloud.ServerType{Name: params.ServerType},
		Image:      &hcloud.Image{ID: imageID},
		Location:   &hcloud.Location{Name: params.Location},
		PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: params.PublicIPv4, EnableIPv6: params.PublicIPv6},
		Labels:     services.AppendManagedLabel(params.GetConfigMapString()),
	}

	if env == "dev" {
		ops.PublicNet.EnableIPv4 = true
		ops.PublicNet.EnableIPv6 = true
		ops.Firewalls = []*hcloud.ServerCreateFirewall{{Firewall: hcloud.Firewall{ID: *params.FirewallID}}}
	}
	server, _, err := hetzner.HClient.Server.Create(ctx, ops)

	if err != nil {
		return err
	}
	return model.SaveDB(ctx, *server.Server, params, DB_ID)

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
