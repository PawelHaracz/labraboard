package dtos

import (
	"github.com/google/uuid"
	"time"
)

type CreateProjectDto struct {
	IacType          int    `json:"type"`
	RepositoryUrl    string `json:"repositoryUrl"`
	RepositoryBranch string `json:"repositoryBranch:"`
	TerraformPath    string `json:"repositoryTerraformPath"`
}

type GetProjectBaseDto struct {
	Id      uuid.UUID `json:"id"`
	IacType int       `json:"type"`
}

type GetProjectDto struct {
	GetProjectBaseDto
	RepositoryUrl    string            `json:"repositoryUrl"`
	RepositoryBranch string            `json:"repositoryBranch:"`
	TerraformPath    string            `json:"repositoryTerraformPath"`
	Envs             map[string]string `json:"envs"`
	Variables        map[string]string `json:"variables"`
}

type SchedulePlan struct {
	When time.Time `json:"when"`
}
