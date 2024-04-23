package dtos

import "time"

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
	RepoPath      string            `json:"repoPath"`
	RepoCommitSha string            `json:"repoCommitSha"`
	Variables     map[string]string `json:"variables"`
	EnvVariables  map[string]string `json:"envVariables"`
}
