package controller

import "time"

type Alert struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	OrgID             int64             `json:"orgId"`
	Alerts            []AlertElement    `json:"alerts"`
	GroupLabels       CommonAnnotations `json:"groupLabels"`
	CommonLabels      CommonLabels      `json:"commonLabels"`
	CommonAnnotations CommonAnnotations `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int64             `json:"truncatedAlerts"`
	Title             string            `json:"title"`
	State             string            `json:"state"`
	Message           string            `json:"message"`
}

type AlertElement struct {
	Status       string      `json:"status"`
	Labels       Labels      `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     time.Time   `json:"startsAt"`
	EndsAt       time.Time   `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Fingerprint  string      `json:"fingerprint"`
	SilenceURL   string      `json:"silenceURL"`
	DashboardURL string      `json:"dashboardURL"`
	PanelURL     string      `json:"panelURL"`
	Values       Values      `json:"values"`
}

type Annotations struct {
	Description string `json:"description"`
	RunbookURL  string `json:"runbook_url"`
	Summary     string `json:"summary"`
}

type Labels struct {
	Alertname string `json:"alertname"`
	Team      string `json:"team"`
	Zone      string `json:"zone"`
}

type Values struct {
	B float64 `json:"B"`
	C int64   `json:"C"`
}

type CommonAnnotations struct {
}

type CommonLabels struct {
	Team string `json:"team"`
}
