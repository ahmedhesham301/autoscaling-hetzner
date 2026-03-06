package services

import (
	"autoscaling-hetzner/grafana"
	"fmt"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"
	"github.com/grafana/grafana-openapi-client-go/client/provisioning"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func SetupAlert(groupId int, monitoringType string, target int16) (string, error) {
	var expression string
	switch monitoringType {
	case "cpu":
		expression = fmt.Sprintf("avg(1-rate(node_cpu_seconds_total{mode=\"idle\",groupId='%d'}[1m]))*100", groupId)
	case "memory":
		expression = fmt.Sprintf("(1-(node_memory_MemAvailable_bytes{groupId=\"%d\"} / node_memory_MemTotal_bytes{groupId=\"%d\"}))*100", groupId, groupId)
	default:
		return "", fmt.Errorf("invalid monitoring type")
	}

	forDur := strfmt.Duration(3 * time.Minute)
	queries := []*models.AlertQuery{
		{
			RefID:         "A",
			DatasourceUID: grafana.PrometheusUid,
			RelativeTimeRange: &models.RelativeTimeRange{
				From: models.Duration(180),
				To:   models.Duration(0),
			},
			Model: map[string]any{
				"datasource": map[string]any{
					"type": "prometheus",
					"uid":  grafana.PrometheusUid,
				},
				"editorMode":    "code",
				"expr":          expression,
				"intervalMs":    1000,
				"legendFormat":  "__auto",
				"maxDataPoints": 43200,
				"range":         true,
				"refId":         "A",
			},
		},
		{
			RefID:         "B",
			DatasourceUID: "__expr__",
			RelativeTimeRange: &models.RelativeTimeRange{
				From: models.Duration(0),
				To:   models.Duration(0),
			},
			Model: map[string]any{
				"conditions": []map[string]any{
					{
						"evaluator": map[string]any{
							"params": []float64{float64(target)},
							"type":   "gt",
						},
						"operator": map[string]any{
							"type": "and",
						},
						"query": map[string]any{
							"params": []string{"A"},
						},
						"reducer": map[string]any{
							"params": []any{},
							"type":   "last",
						},
						"type": "query",
					},
				},
				"datasource": map[string]any{
					"name": "Expression",
					"type": "__expr__",
					"uid":  "__expr__",
				},
				"expression":    "A",
				"intervalMs":    1000,
				"maxDataPoints": 43200,
				"refId":         "B",
				"type":          "classic_conditions",
			},
		},
	}

	rule := &models.ProvisionedAlertRule{
		Title:                conv.Pointer(fmt.Sprintf("%s usage", monitoringType)),
		Condition:            conv.Pointer("B"),
		Data:                 queries,
		ExecErrState:         conv.Pointer("Alerting"),
		NoDataState:          conv.Pointer("NoData"),
		FolderUID:            &grafana.FolderUid,
		For:                  &forDur,
		OrgID:                conv.Pointer(int64(1)),
		RuleGroup:            conv.Pointer("default"),
		NotificationSettings: &models.AlertRuleNotificationSettings{Receiver: conv.Pointer("server")},
		UID:                  strconv.Itoa(groupId),
		Labels:               map[string]string{"groupId": strconv.Itoa(groupId)},
	}
	err := rule.Validate(strfmt.Default)
	if err != nil {
		return "", err
	}

	resp, err := grafana.GClient.Provisioning.PostAlertRule(provisioning.NewPostAlertRuleParams().WithBody(rule))
	if err != nil {
		return "", err
	}
	return resp.Payload.UID, nil
}
