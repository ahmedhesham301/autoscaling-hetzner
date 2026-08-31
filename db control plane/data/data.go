package data

import "fmt"

type ServiceInfo struct {
	AppName     string            `json:"app_name"`
	AppVersions []string          `json:"app_versions"`
	Vars        map[string]string `json:"vars"`
}

type CreateServiceParams struct {
	AppName         string            `json:"app_name"`
	AppVersion      string            `json:"app_version"`
	Location        string            `json:"location"`
	ServerType      string            `json:"server_type"`
	PublicIPv4      bool              `json:"public_ipv4"`
	PublicIPv6      bool              `json:"public_ipv6"`
	NodeExporter    bool              `json:"node_exporter"`
	ServiceExporter bool              `json:"service_exporter"`
	ExtraLabels     map[string]string `json:"extra_labels"`
}

func (params CreateServiceParams) GetConfigMap() map[string]any {
	return map[string]any{
		"app_name":         params.AppName,
		"app_version":      params.AppVersion,
		"node_exporter":    params.NodeExporter,
		"service_exporter": params.ServiceExporter,
	}
}

func (params CreateServiceParams) GetConfigMapString() map[string]string {
	result := make(map[string]string)
	for k, v := range params.GetConfigMap() {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
