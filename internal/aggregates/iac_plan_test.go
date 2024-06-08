package aggregates

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"labraboard/internal/entities"
	"labraboard/internal/helpers"
	"labraboard/internal/models"
	"os"
	"path"
	"runtime"
	"testing"
)

var planJson []byte

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	pathTestPlan := path.Join(dir, "testingArtifacts", "terraform_plan.json")
	planJson, _ = os.ReadFile(pathTestPlan)
}

func TestIaCPlanAddPlan(t *testing.T) {
	//arrange
	aggregate, _ := NewIacPlan(uuid.New(), Terraform, nil)

	//act
	aggregate.AddPlan(planJson, nil)

	assert.Equal(t, aggregate.planJson, planJson)
}

func TestIaCPlan_ChangesBasedOnPlan_ValidChangeCount(t *testing.T) {
	//arrange
	serializer := helpers.NewSerializer[entities.IacTerraformOutputJson]()
	aggregate, _ := NewIacPlan(uuid.New(), Terraform, nil)
	buf := bytes.NewBuffer(planJson)
	plan, err := serializer.DeserializeJsonl(buf)
	if err != nil {
		t.Error(err)
	}
	planChanges := models.NewIacTerraformPlanJson(planJson, plan, nil)

	//act
	aggregate.AddChanges(planChanges.GetChanges()...)

	//assert
	add, change, del := aggregate.GetChanges()
	assert.Equal(t, add, 5)
	assert.Equal(t, change, 0)
	assert.Equal(t, del, 0)
}

func TestIaCPlan_ChangesBasedOnPlan_ValidResourceChanges(t *testing.T) {
	//arrange
	serializer := helpers.NewSerializer[entities.IacTerraformOutputJson]()
	aggregate, _ := NewIacPlan(uuid.New(), Terraform, nil)
	buf := bytes.NewBuffer(planJson)
	plan, err := serializer.DeserializeJsonl(buf)
	if err != nil {
		t.Error(err)
	}
	planChanges := models.NewIacTerraformPlanJson(planJson, plan, nil)

	//act
	aggregate.AddChanges(planChanges.GetChanges()...)

	//assert
	assert.Equal(t, len(aggregate.changes), 5)
	assert.Equal(t, aggregate.changes[0].ResourceType, "azurerm_resource_group")
	assert.Equal(t, aggregate.changes[1].ResourceType, "azurerm_virtual_network")
	assert.Equal(t, aggregate.changes[2].ResourceType, "azurerm_subnet")
	assert.Equal(t, aggregate.changes[3].ResourceType, "azurerm_network_interface")
	assert.Equal(t, aggregate.changes[4].ResourceType, "azurerm_linux_virtual_machine")
}
