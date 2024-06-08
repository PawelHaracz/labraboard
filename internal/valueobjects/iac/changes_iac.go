package iac

type PlanTypeAction string

var (
	Create PlanTypeAction = "create"
	Update PlanTypeAction = "update"
	Delete PlanTypeAction = "delete"
)

type ChangesIac struct {
	ResourceType string
	ResourceName string
	Provider     string
	Action       PlanTypeAction
}
