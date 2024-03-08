package iac

import "github.com/google/uuid"

type IacType string

const (
	Tofu      IacType = "tofu"
	Terraform IacType = "terraform"
)

type Plan struct {
	Type IacType
	Id   uuid.UUID
	plan interface{}
}

type LabraboardIacService interface {
	Plan(planId uuid.UUID) (*Plan, error)
}
