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
