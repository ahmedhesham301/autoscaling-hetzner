package temporal

import (
	"log/slog"
	"time"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"go.temporal.io/sdk/workflow"
)

func CreateServiceWorkflow(ctx workflow.Context, params data.CreateServiceParams) error {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Check if image exists
	var imageID *int64
	err := workflow.ExecuteActivity(ctx, checkImageExist, params).Get(ctx, &imageID)
	if err != nil {
		slog.Error("error checking if image exists", "err", err)
		return err
	}
	// If not build it
	if imageID == nil {
		err := workflow.ExecuteActivity(ctx, buildImage, params).Get(ctx, &imageID)
		if err != nil {
			slog.Error("error building image", "err", err)
			return err
		}
	}
	// Deploy it
	err = workflow.ExecuteActivity(ctx, deployService, params, imageID).Get(ctx, nil)
	if err != nil {
		slog.Error("error building image", "err", err)
		return err
	}
	return nil
}
