package main

import (
	"fmt"
	"testing"
)

func TestParseThing(t *testing.T) {
	tests := []struct {
		name     string
		val      []byte
		expected string
	}{
		{name: "Positive Integer", val: []byte("i42e"), expected: "integer"},
		{name: "Zero", val: []byte("i0e"), expected: "integer"},
		{name: "Byte String", val: []byte("4:spam"), expected: "byte_string"},
		{name: "Byte String", val: []byte("5:thing"), expected: "byte_string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseThing(tt.val)

			fmt.Println(v, err)

			if err != nil {
				t.Errorf("ParseString(%v) = %v\n", tt.val, v)
			}
		})
	}
}
