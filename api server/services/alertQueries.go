package services

import (
	"autoscaling-hetzner/grafana"
	"autoscaling-hetzner/model"
	"fmt"

	"github.com/grafana/grafana-openapi-client-go/models"
)

func setupAlertQuery(group *model.Group) ([]*models.AlertQuery, error) {
	var expression string
	switch group.MonitoringType {
	case "cpu":
		expression = fmt.Sprintf("avg(1-rate(node_cpu_seconds_total{mode=\"idle\",groupId='%d'}[1m]))*100", group.Id)
	case "memory":
		expression = fmt.Sprintf("(1-(node_memory_MemAvailable_bytes{groupId=\"%d\"} / node_memory_MemTotal_bytes{groupId=\"%d\"}))*100", group.Id, group.Id)
	default:
		return nil, fmt.Errorf("invalid monitoring type")
	}

	var condition map[string]any
	switch group.ScalingAlgorithm {
	case "simple":
		condition = map[string]any{
			"evaluator": map[string]any{
				"params": []int16{*group.ScaleDownThreshold, *group.ScaleUpThreshold},
				"type":   "outside_range",
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
		}
	default:
		return nil, fmt.Errorf("invalid scaling algorithm")
	}

	queries := []*models.AlertQuery{
		{
			RefID:         "A",
			DatasourceUID: grafana.PrometheusUid,
			RelativeTimeRange: &models.RelativeTimeRange{
				From: models.Duration(1),
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
				"conditions": []map[string]any{condition},
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
	return queries, nil
}
