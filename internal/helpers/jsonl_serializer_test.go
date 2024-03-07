package helpers

import (
	"labraboard/internal/entities"
	"os"
	"testing"
)

func TestIacTerraformPlanJsons(t *testing.T) {
	f, err := os.Open("/Users/pawelharacz/src/labraboard/testingArtifacts/terraform_plan.json")
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	defer f.Close()

	serializer := NewSerializer[entities.IacTerraformPlanJson]()
	iacTerraformPlanJsons, err := serializer.SerializeJsonl(f)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(iacTerraformPlanJsons) != 7 {
		t.Fail()
	}

}
