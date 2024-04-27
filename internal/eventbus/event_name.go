package eventbus

type EventName string

const (
	TRIGGERED_PLAN  EventName = "triggered_plan"
	LEASE_PLAN_LOCK EventName = "lease_plan_lock"
)
