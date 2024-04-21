package models

import "labraboard/internal/entities"

type IacTerraformPlanJson struct {
	plan    []byte
	changes []entities.IacTerraformPlanJson
}

func NewIacTerraformPlanJson(plan []byte, changes []entities.IacTerraformPlanJson) *IacTerraformPlanJson {
	return &IacTerraformPlanJson{
		plan:    plan,
		changes: changes,
	}
}

func (p *IacTerraformPlanJson) GetPlan() []byte {
	return p.plan
}
func (p *IacTerraformPlanJson) GetChanges() []entities.IacTerraformPlanJson {
	var value = make([]entities.IacTerraformPlanJson, len(p.changes))
	for i, v := range p.changes {
		value[i] = v
	}
	return value
}
