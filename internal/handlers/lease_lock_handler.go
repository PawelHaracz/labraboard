package handlers

import (
	eb "labraboard/internal/eventbus"
	"labraboard/internal/models"
	"labraboard/internal/repositories"
)

//todo implement generic one handler with object instead multiple
import (
	"context"
	"encoding/json"
	"fmt"
	"labraboard/internal/eventbus/events"
)

type terraformStateLeaseLockHandler struct {
	eventSubscriber eb.EventSubscriber
	unitOfWork      *repositories.UnitOfWork
}

func newTerraformStateLeaseLockHandler(eventSubscriber eb.EventSubscriber, unitOfWork *repositories.UnitOfWork) (*terraformStateLeaseLockHandler, error) {
	return &terraformStateLeaseLockHandler{
		eventSubscriber,
		unitOfWork,
	}, nil
}

func (handler *terraformStateLeaseLockHandler) Handle(ctx context.Context) {
	locks := handler.eventSubscriber.Subscribe(events.LEASE_LOCK, ctx)
	for msg := range locks {
		var event = events.LeasedLock{}
		err := json.Unmarshal(msg, &event)
		if err != nil {
			panic(fmt.Errorf("cannot handle message type %T", event))
		}
		fmt.Println("Received message:", msg)
		go handler.handle(event)
	}
}

func (handler *terraformStateLeaseLockHandler) handle(event events.LeasedLock) {
	if event.Type != models.Terraform {
		return
	}

	item, err := handler.unitOfWork.TerraformStateDbRepository.Get(event.Id)
	if err != nil {
		panic(err) //todo logging
	}

	info, err := item.GetLockInfo()
	if err != nil {
		//todo logging
		return
	}

	if info == nil {
		return
	}

	if err = item.SetLockInfo(nil); err != nil {
		//todo logging
		return
	}
}
