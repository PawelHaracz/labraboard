package events

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type EventName string

type Event interface {
	MarshalBinary() ([]byte, error)
}

func Unmarshal(name EventName, b []byte) (Event, error) {
	switch name {
	case LEASE_LOCK:
		var l LeasedLock
		err := json.Unmarshal(b, &l)
		return l, err
	case SCHEDULED_PLAN:
		var l ScheduledPlan
		err := json.Unmarshal(b, &l)
		return l, err
	case TRIGGERED_PLAN:
		var l PlanTriggered
		err := json.Unmarshal(b, &l)
		return l, err
	case IAC_APPLY_SCHEDULED:
		var l IacApplyScheduled
		err := json.Unmarshal(b, &l)
		return l, err
	default:
		return nil, errors.New(fmt.Sprintf("%s is not supported", name))
	}
}
