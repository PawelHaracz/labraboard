package helpers

import (
	"labraboard/internal/entities"
	"os"
	"testing"
)

func TestIacTerraformPlanJsons(t *testing.T) {
	f, err := os.Open("../../testingArtifacts/terraform_plan.json")
	if err != nil {
		t.Errorf("failed to open file: %v", err)
		t.Fail()
	}
	defer f.Close()

	serializer := NewSerializer[entities.IacTerraformOutputJson]()
	iacTerraformPlanJsons, err := serializer.DeserializeJsonl(f)
	if err != nil {
		t.Errorf("failed to deserialize: %v", err)
	}

	if len(iacTerraformPlanJsons) != 7 {
		t.Fail()
	}
}

func TestIacTerraformApplyErrorJsons(t *testing.T) {
	f, err := os.Open("../../testingArtifacts/terraform_apply_error.json")
	if err != nil {
		t.Errorf("failed to open file: %v", err)
		t.Fail()
	}
	defer f.Close()

	serializer := NewSerializer[entities.IacTerraformDiagnosticJson]()
	iacTerraformPlanJsons, err := serializer.DeserializeJsonl(f)
	if err != nil {
		t.Errorf("failed to deserialize: %v", err)
	}

	if len(iacTerraformPlanJsons) != 1 {
		t.Fail()
	}
}
