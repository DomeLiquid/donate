package json

import (
	"github.com/bytedance/sonic"
)

func Marshal[T any](j T) ([]byte, error) {
	return sonic.Marshal(j)
}

func Unmarshal[T any](b []byte, data T) error {
	return sonic.Unmarshal(b, &data)
}

