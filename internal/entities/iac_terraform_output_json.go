package entities

import (
	"time"
)

type IacPlatformType string

var (
	Version       IacPlatformType
	PlannedChange IacPlatformType
	ChangeSummary IacPlatformType
)

type IacTerraformOutputJson struct {
	Message        string                             `json:"@message"`
	Module         string                             `json:"@module"`
	Timestamp      time.Time                          `json:"@timestamp"`
	Type           IacPlatformType                    `json:"type"`
	Change         IacTerraformChangeJson             `json:"change"`
	SummaryChanges IacTerraformSummaryChangesJson     `json:"changes"`
	Outputs        map[string]IacTerraformOutputValue `json:"outputs"`
}

type IacTerraformChangeJson struct {
	Resource IacTerraformChangeResourceJson `json:"resource"`
	Action   string                         `json:"action"`
}

type IacTerraformChangeResourceJson struct {
	Addr         string `json:"addr"`
	Module       string `json:"module"`
	Resource     string `json:"resource"`
	Provider     string `json:"implied_provider"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	ResourceKey  string `json:"resource_key"`
}

type IacTerraformSummaryChangesJson struct {
	Add       int    `json:"add"`
	Change    int    `json:"change"`
	Remove    int    `json:"remove"`
	Operation string `json:"operation"`
}

type IacTerraformOutputValue struct {
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}
