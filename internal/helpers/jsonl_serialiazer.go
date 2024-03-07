package helpers

import (
	"encoding/json"
	"io"
)

type Serializer[T any] struct {
}

func NewSerializer[T any]() *Serializer[T] {
	return &Serializer[T]{}
}

func (*Serializer[T]) SerializeJsonl(reader io.Reader) ([]T, error) {
	var jsons []T
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		var item T
		if err := decoder.Decode(&item); err != nil {
			return nil, err
		}
		jsons = append(jsons, item)
	}
	return jsons, nil
}
