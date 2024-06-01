package models

import (
	"labraboard/internal/entities"
)

type IacTerraformPlanJson struct {
	plan    []byte
	changes []entities.IacTerraformPlanJson
	planRaw []byte
}

func NewIacTerraformPlanJson(plan []byte, changes []entities.IacTerraformPlanJson, planRaw []byte) *IacTerraformPlanJson {
	return &IacTerraformPlanJson{
		plan:    plan,
		changes: changes,
		planRaw: planRaw,
	}
}

func (p *IacTerraformPlanJson) GetPlan() (planJson []byte, planRaw []byte) {
	return p.plan, p.planRaw
}
func (p *IacTerraformPlanJson) GetChanges() []entities.IacTerraformPlanJson {
	var value = make([]entities.IacTerraformPlanJson, len(p.changes))
	for i, v := range p.changes {
		value[i] = v
	}
	return value
}
