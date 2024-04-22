package aggregates

import (
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	vo "labraboard/internal/valueobjects"
	"testing"
)

func TestNewIaC_CreateEmptyProject_ShouldGenerate(t *testing.T) {
	projectId := uuid.New()

	iac, err := NewIac(projectId, vo.Terraform, make([]*vo.Plans, 0), make([]*vo.IaCEnv, 0), nil, make([]*vo.IaCVariable, 0))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, iac.id, projectId)
}

func TestNewIaC_AddRepository_WhenRepoNotDefined(t *testing.T) {
	projectId := uuid.New()

	iac, err := NewIac(projectId, vo.Tofu, nil, nil, nil, nil)
	err = iac.AddRepo("https://github.com/alfonsof/terraform-azure-examples", "master", "code/03-one-webserver")
	url, branch, path := iac.GetRepo()

	assert.NotEqual(t, err, nil)
	assert.NotEqual(t, iac, nil)
	assert.Equal(t, url, "https://github.com/alfonsof/terraform-azure-examples")
	assert.Equal(t, branch, "master")
	assert.Equal(t, path, "code/03-one-webserver")

}
