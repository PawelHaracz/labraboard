package handlers

import (
	"github.com/stretchr/testify/assert"
	"labraboard/internal/eventbus/events"
	"testing"
)

func TestEventHandlerFactory(t *testing.T) {
	var factory = NewEventHandlerFactory(nil, nil, nil)
	t.Run("Register Event that doesn't exist, should return null", func(t *testing.T) {
		var Z events.EventName = "Z"
		act, err := factory.RegisterHandler(Z)

		assert.Equal(t, nil, act)
		assert.Error(t, MissingHandlerImplementedFactory, err.Error())
	})
	t.Run("Register Event that exist, should return handler", func(t *testing.T) {
		act, err := factory.RegisterHandler(events.SCHEDULED_PLAN)

		assert.Equal(t, nil, err)
		assert.NotEqual(t, act, nil)
	})

	t.Run("Register Twice the same event, should raise the error", func(t *testing.T) {
		act, err := factory.RegisterHandler(events.LEASE_LOCK)

		assert.Equal(t, nil, err)
		assert.NotEqual(t, act, nil)

		act, err = factory.RegisterHandler(events.LEASE_LOCK)
		assert.Equal(t, nil, act)
		assert.Error(t, HandlerAlreadyRegistered, err.Error())
	})

	t.Run("Register all handlers, should be ok", func(t *testing.T) {
		var expected = len(factory.allowedEvents)
		act, err := factory.RegisterAllHandlers()
		assert.Nil(t, err)
		assert.Equal(t, expected, len(act))

	})
}
