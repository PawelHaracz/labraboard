package managers

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"labraboard/internal/eventbus/events"
	"labraboard/internal/models"
	"testing"
	"time"
)

func TestJsonUnmarshal(t *testing.T) {
	layout := "2006-01-02T15:04:05.999999Z"
	createdTime, err := time.Parse(layout, "2024-02-05T20:04:43.120857Z")
	var lease = events.LeasedLock{
		Id:        uuid.New(),
		Type:      models.Terraform,
		LeaseTime: createdTime,
	}
	b, err := lease.MarshalBinary()
	if err != nil {
		t.Fatal(err.Error())
	}
	var ta = task{
		EventName: events.LEASE_LOCK,
		Content:   b,
		WaitTime:  10,
	}

	z, err := json.Marshal(ta)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = json.Unmarshal(z, &ta)
	if err != nil {
		t.Fatal(err.Error())
	}

	l, err := events.Unmarshal(ta.EventName, ta.Content)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, ta.EventName, events.LEASE_LOCK)
	assert.Equal(t, (l.(events.LeasedLock)).LeaseTime, lease.LeaseTime)
	assert.Equal(t, (l.(events.LeasedLock)).Id, lease.Id)
	assert.Equal(t, (l.(events.LeasedLock)).Type, lease.Type)

}
