package temporal

import (
	"time"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"go.temporal.io/sdk/workflow"
)

func CreateServiceWorkflow(ctx workflow.Context, params data.CreateDBParams, DB_ID int, env string, templatesPath string) error {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)

	// Check if image exists
	var imageID *int64
	err := workflow.ExecuteActivity(ctx, checkImageExist, params).Get(ctx, &imageID)
	if err != nil {
		logger.Error("error checking if image exists", "err", err)
		return err
	}
	// If not build it
	if imageID == nil {
		err := workflow.ExecuteActivity(ctx, buildImage, params, env, templatesPath).Get(ctx, &imageID)
		if err != nil {
			logger.Error("error building image", "err", err)
			return err
		}
	}
	// create a firewall that allows traffic if env is dev
	if env == "dev" {
		err = workflow.ExecuteActivity(ctx, GetOrCreateAllowAllFirewall).Get(ctx, &params.FirewallID)
		if err != nil {
			logger.Error("error building image", "err", err)
			return err
		}
	}
	// Deploy it
	err = workflow.ExecuteActivity(ctx, deployDB, params, imageID, DB_ID).Get(ctx, nil)
	if err != nil {
		logger.Error("error building image", "err", err)
		return err
	}
	return nil
}
