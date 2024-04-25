package dtos

import (
	"labraboard/internal/models"
	"time"
)

type PlanDto struct {
	Id        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedOn time.Time `json:"createdOn"`
}

type PlanWithOutputDto struct {
	Id        string      `json:"id"`
	Status    string      `json:"status"`
	CreatedOn time.Time   `json:"createdOn"`
	Outputs   interface{} `json:"outputs"`
}

type CreatePlan struct {
	RepoPath       string            `json:"repoPath"`
	RepoCommit     string            `json:"repoCommit"`
	RepoCommitType models.CommitType `json:"repoCommitType"`
	Variables      map[string]string `json:"variables"`
}
