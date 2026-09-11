package main

import (
	"fmt"
	"testing"
)

//func TestParseType(t *testing.T) {
//	tests := []struct {
//		name     string
//		val      []byte
//		expected BencodeType
//	}{
//		{name: "Positive Integer", val: []byte("i42e"), expected: EncodedInt},
//		{name: "Zero", val: []byte("i0e"), expected: EncodedInt},
//		{name: "Byte String", val: []byte("4:spam"), expected: EncodedByteString},
//		{name: "Byte String", val: []byte("5:thing"), expected: EncodedByteString},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			v, err := ParseType(tt.val)
//
//			if err != nil {
//				t.Errorf("ParseString(%v) = %v\n", tt.val, v)
//			}
//		})
//	}
//}

func TestDecoder(t *testing.T) {
	d := Decoder{
		buf: make([]byte, 0),
		pos: 10,
	}

	d.buf = append(d.buf, []byte("17:abcdefghijklmnopqrstuvwxyz")...)

	fmt.Println("Buffer is: ", d.buf, string(d.buf))

	_, err := d.ParseByteString()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
