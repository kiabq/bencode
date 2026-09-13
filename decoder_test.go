package main

import (
	"fmt"
	"testing"
)

func TestDecoder(t *testing.T) {
	t.Run("ParseInt", func(t *testing.T) {
		tests := []struct {
			name     string
			val      []byte
			pos      int
			fail     bool
			expected any
		}{
			{name: "positive integer", val: []byte("i105e"), fail: false},
			{name: "zero integer", val: []byte("i0e"), fail: false},
			{name: "in a list", val: []byte("i42e"), pos: 3, fail: false},
			{name: "leading zero integer", val: []byte("i042e"), fail: true},
			{name: "max int overflow", val: []byte("i9223372036854775808e"), fail: true},
			{name: "contains non-digit characters", val: []byte("i33892d21e"), fail: true},
			{name: "negative zero", val: []byte("i-0e"), fail: true},
			{name: "double negative", val: []byte("i--2e"), fail: true},
			{name: "contains minus", val: []byte("i4-2e"), fail: true},
			{name: "missing terminator character", val: []byte("i21321312"), fail: true},
			{name: "incorrect terminator", val: []byte("i42d"), fail: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := Decoder{
					buf: make([]byte, tt.pos),
					pos: tt.pos,
				}

				d.buf = append(d.buf, tt.val...)
				integer, err := d.ParseInt()
				fmt.Println("position: ", d.pos)
				fmt.Println("buffer @ position: ", string(d.buf[:d.pos]))
				fmt.Println("integer: ", integer)
				fmt.Println("error: ", err)
				if err != nil && !tt.fail {
					t.Errorf("ParseInt: %v\n", err)
				} else {
					t.Logf("ParseInt: %d\n", integer)
				}
			})
		}
	})
}
