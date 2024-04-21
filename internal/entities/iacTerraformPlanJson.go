package entities

import (
	"time"
)

type IacPlatformPlanType string

var (
	Version       IacPlatformPlanType
	PlannedChange IacPlatformPlanType
	ChangeSummary IacPlatformPlanType
)

type IacTerraformPlanJson struct {
	Message        string                             `json:"@message"`
	Module         string                             `json:"@module"`
	Timestamp      time.Time                          `json:"@timestamp"`
	Type           IacPlatformPlanType                `json:"type"`
	Change         IacTerraformPlanChangeJson         `json:"change"`
	SummaryChanges IacTerraformPlanSummaryChangesJson `json:"changes"`
}

type IacTerraformPlanChangeJson struct {
	Resource IacTerraformPlanChangeResourceJson `json:"resource"`
	Action   string                             `json:"action"`
}

type IacTerraformPlanChangeResourceJson struct {
	Addr         string `json:"addr"`
	Module       string `json:"module"`
	Resource     string `json:"resource"`
	Provider     string `json:"implied_provider"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	ResourceKey  string `json:"resource_key"`
}

type IacTerraformPlanSummaryChangesJson struct {
	Add       int    `json:"add"`
	Change    int    `json:"change"`
	Remove    int    `json:"remove"`
	Operation string `json:"operation"`
}
