package config

type ConfigModel struct {
	ServiceName string `json:"service_name"`
	Env         string `json:"env"`
	Config      map[string]any `json:"config"`
}