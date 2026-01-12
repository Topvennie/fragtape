package queue

import "encoding/json"

type Codec[T any] interface {
	Marshal(T) ([]byte, error)
	Unmarshal([]byte) (T, error)
}

type JSONCodec[T any] struct{}

// Interface compliance
var _ Codec[any] = (*JSONCodec[any])(nil)

func (JSONCodec[T]) Marshal(t T) ([]byte, error) {
	return json.Marshal(t)
}

func (JSONCodec[T]) Unmarshal(data []byte) (T, error) {
	var t T
	err := json.Unmarshal(data, &t)

	return t, err
}
