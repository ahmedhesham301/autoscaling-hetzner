package services

import (
	"autoscaling-hetzner/grafana"
	"autoscaling-hetzner/model"
	"fmt"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"
	"github.com/grafana/grafana-openapi-client-go/client/provisioning"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func SetupAlert(group *model.Group) (string, error) {

	queries, err := setupAlertQuery(group)
	if err != nil {
		return "", err
	}

	forDur := strfmt.Duration(3 * time.Minute)

	rule := &models.ProvisionedAlertRule{
		Title:                conv.Pointer(fmt.Sprintf("%s usage", group.MonitoringType)),
		Condition:            conv.Pointer("B"),
		Data:                 queries,
		ExecErrState:         conv.Pointer("Alerting"),
		NoDataState:          conv.Pointer("NoData"),
		FolderUID:            &grafana.FolderUid,
		For:                  &forDur,
		OrgID:                conv.Pointer(int64(1)),
		RuleGroup:            conv.Pointer("default"),
		NotificationSettings: &models.AlertRuleNotificationSettings{Receiver: conv.Pointer("server")},
		UID:                  strconv.Itoa(group.Id),
		Labels:               map[string]string{"groupId": strconv.Itoa(group.Id)},
	}
	err = rule.Validate(strfmt.Default)
	if err != nil {
		return "", err
	}

	resp, err := grafana.GClient.Provisioning.PostAlertRule(provisioning.NewPostAlertRuleParams().WithBody(rule))
	if err != nil {
		return "", err
	}
	return resp.Payload.UID, nil
}
