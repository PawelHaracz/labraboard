package entities

import (
	"encoding/json"
	"io"
	"time"
)

type IacTerraformPlanJson struct {
	Message        string                              `json:"@message"`
	Module         string                              `json:"@module"`
	Timestamp      time.Time                           `json:"@timestamp"`
	Type           string                              `json:"type"`
	Change         *IacTerraformPlanChangeJson         `json:"change"`
	SummaryChanges *IacTerraformPlanSummaryChangesJson `json:"changes"`
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

func SerializeIacTerraformPlanJsons(file io.Reader) ([]IacTerraformPlanJson, error) {
	var iacTerraformPlanJsons []IacTerraformPlanJson
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var iacTerraformPlanJson IacTerraformPlanJson
		if err := decoder.Decode(&iacTerraformPlanJson); err != nil {
			return nil, err
		}
		iacTerraformPlanJsons = append(iacTerraformPlanJsons, iacTerraformPlanJson)
	}
	return iacTerraformPlanJsons, nil
}
