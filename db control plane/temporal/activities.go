package temporal

import (
	"context"
	"fmt"
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

func buildImage(ctx context.Context, params data.CreateDBParams) (*int64, error) {
	logger := activity.GetLogger(ctx)
	templatesPath := os.Getenv("PACKER_TEMPLATES_PATH")

	cmd := exec.CommandContext(ctx, "packer", "build", "-machine-readable",
		"-var", fmt.Sprintf("config=%v", utils.ConvertMapToJsonString(params.GetConfigMap())),
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

	if env == "dev" {
		params.PublicIPv4 = true
		params.PublicIPv6 = true
	}

	server, _, err := hetzner.HClient.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       random.AddRandomLetters(params.AppName + "-" + params.AppVersion),
		ServerType: &hcloud.ServerType{Name: params.ServerType},
		Image:      &hcloud.Image{ID: imageID},
		Location:   &hcloud.Location{Name: params.Location},
		PublicNet:  &hcloud.ServerCreatePublicNet{EnableIPv4: params.PublicIPv4, EnableIPv6: params.PublicIPv6},
		Labels:     services.AppendManagedLabel(params.GetConfigMapString()),
	})

	if err != nil {
		return err
	}
	return model.SaveDB(ctx, *server.Server, params, DB_ID)

}
