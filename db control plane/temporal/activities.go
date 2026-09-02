package temporal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/utils"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)


func checkImageExist(ctx context.Context, params data.CreateServiceParams) (*int64, error) {
	fmt.Println(services.ConvertToHetznerLabels(params.GetConfigMap()))
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

func buildImage(ctx context.Context, params data.CreateServiceParams) (*int64, error) {
	templatesPath, exists := os.LookupEnv("PACKER_TEMPLATES_PATH")
	if !exists {
		slog.Error("env var PACKER_TEMPLATES_PATH is not set")
		os.Exit(1)
	}

	cmd := exec.Command("packer", "build",
		"-var", fmt.Sprintf("config=%v", utils.ConvertMapToJsonString(params.GetConfigMap())),
		templatesPath+"/"+params.AppName+"/main.pkr.hcl")

	output, err := cmd.Output()
	if err != nil {
		slog.Error("Command failed", "err", err)
	}

	fmt.Println(string(output))

	var id int64 = 1
	return &id, nil
}

func deployService(ctx context.Context, params data.CreateServiceParams, imageID int64) error {
	_, _, err := hetzner.HClient.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       params.AppName + params.AppVersion,
		ServerType: &hcloud.ServerType{Name: params.ServerType},
		Image:      &hcloud.Image{ID: imageID},
		Location:   &hcloud.Location{Name: params.Location},
		PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: params.PublicIPv4, EnableIPv6: params.PublicIPv6},
		Labels:     services.AppendManagedLabel(params.GetConfigMapString()),
	})
	return err
}
