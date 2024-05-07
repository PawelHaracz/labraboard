package managers

import (
	"golang.org/x/net/context"
	"labraboard/internal/eventbus/events"
	"time"
)

type memoryDelayTask struct {
}

func NewMemoryDelayTask() DelayTaskManager {
	return memoryDelayTask{}
}

func (m memoryDelayTask) Listen(ctx context.Context) {

}

func (m memoryDelayTask) Publish(EventName events.EventName, Content events.Event, WaitTime time.Duration, ctx context.Context) {

}
