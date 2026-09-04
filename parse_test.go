package main

import (
	"fmt"
	"testing"
)

func TestParseType(t *testing.T) {
	tests := []struct {
		name     string
		val      []byte
		expected BencodeType
	}{
		{name: "Positive Integer", val: []byte("i42e"), expected: EncodedInt},
		{name: "Zero", val: []byte("i0e"), expected: EncodedInt},
		{name: "Byte String", val: []byte("4:spam"), expected: EncodedByteString},
		{name: "Byte String", val: []byte("5:thing"), expected: EncodedByteString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseType(tt.val)

			fmt.Println(v, err)

			if err != nil {
				t.Errorf("ParseString(%v) = %v\n", tt.val, v)
			}
		})
	}
}
