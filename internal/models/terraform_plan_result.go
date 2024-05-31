package models

import (
	"bufio"
	"github.com/pkg/errors"
	"labraboard/internal/entities"
	"os"
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

func (p *IacTerraformPlanJson) SavePlanAsTfPlan(path string) error {
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err = fo.Close(); err != nil {
			err = errors.Wrap(err, "problem with close file")
		}
	}()

	w := bufio.NewWriter(fo)
	if _, err = w.Write(p.planRaw); err != nil {
		err = errors.Wrap(err, "problem with write file")
	}

	return err
}
