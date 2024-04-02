package dtos

import (
	"github.com/google/uuid"
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
	RepositoryUrl    string                   `json:"repositoryUrl"`
	RepositoryBranch string                   `json:"repositoryBranch:"`
	TerraformPath    string                   `json:"repositoryTerraformPath"`
	Envs             []map[string]interface{} `json:"envs"`
}
