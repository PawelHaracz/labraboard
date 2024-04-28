package events

type EventName string

type Event interface {
	MarshalBinary() ([]byte, error)
}
