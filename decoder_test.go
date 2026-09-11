package main

import (
	"errors"
	"math"
	"strconv"
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
		pos: 0,
	}

	d.buf = append(d.buf, []byte("l4:spam4:eggse")...)

	_, err := d.ParseList()
	if err != nil {
		t.Errorf("ParseByteString: %v\n", err)
	}
}

func TestInt64MaxValue(t *testing.T) {
	//var maxInt64 string = "9223372036854775807"
	var overflowInt64 = strconv.Itoa(math.MaxInt - 1)
	_, err := strconv.Atoi("9223372036854775808")
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			t.Errorf("overflowString -> int64 exceeds range")
		} else {
			t.Errorf("unknown error: %v\n", err)
		}
	}

	_, err = strconv.ParseInt(overflowInt64, 10, 64)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			t.Errorf("int64 exceeds range")
		}
	}
}
