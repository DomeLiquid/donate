package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{struct{ Name string }{"John"}, `{"Name":"John"}`},
		{struct{ Age int }{25}, `{"Age":25}`},
		{struct{ Score float64 }{3.14}, `{"Score":3.14}`},
	}

	for _, test := range tests {
		result, err := Marshal(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}

func TestUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int64  `json:"age"`
	}
	tests := []struct {
		input    string
		expected TestStruct
	}{
		{`{"Name":"John", "age": 12}`, TestStruct{
			Name: "John",
			Age:  12,
		}},
	}

	for _, test := range tests {
		var result TestStruct
		err := Unmarshal([]byte(test.input), &result)
		assert.NoError(t, err)
		assert.Equal(t, test.expected.Age, result.Age)
		assert.Equal(t, test.expected.Name, result.Name)
	}
}
