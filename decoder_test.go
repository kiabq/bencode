package main

import (
	"fmt"
	"testing"
)

type DecoderTest struct {
	name string
	val  []byte
	pos  int
	fail bool
}

func TestDecoder(t *testing.T) {
	t.Run("ParseInt", func(t *testing.T) {
		tests := []DecoderTest{
			// valid
			{name: "positive integer", val: []byte("i105e"), fail: false},
			{name: "zero integer", val: []byte("i0e"), fail: false},
			{name: "in a list", val: []byte("i42e"), pos: 3, fail: false},
			{name: "negative integer", val: []byte("i-42e"), fail: false},
			{name: "max int64", val: []byte("i9223372036854775807e"), fail: false},
			// invalid
			{name: "leading zero integer", val: []byte("i042e"), fail: true},
			{name: "max int overflow", val: []byte("i9223372036854775808e"), fail: true},
			{name: "min int underflow", val: []byte("i-9223372036854775809e"), fail: true},
			{name: "contains non-digit characters", val: []byte("i33892d21e"), fail: true},
			{name: "negative zero", val: []byte("i-0e"), fail: true},
			{name: "double negative", val: []byte("i--2e"), fail: true},
			{name: "contains minus", val: []byte("i4-2e"), fail: true},
			{name: "plus prefix", val: []byte("i+42e"), fail: true},
			{name: "minus no digits", val: []byte("i-e"), fail: true},
			{name: "leading zero after minus", val: []byte("i-042e"), fail: true},
			{name: "missing terminator character", val: []byte("i21321312"), fail: true},
			{name: "incorrect terminator", val: []byte("i42d"), fail: true},
			{name: "no integer", val: []byte("ie"), fail: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := Decoder{
					buf: make([]byte, tt.pos),
					pos: tt.pos,
				}

				d.buf = append(d.buf, tt.val...)

				integer, err := d.ParseInt()
				if err != nil && !tt.fail {
					t.Errorf("ParseInt: %v\n", err)
				} else {
					t.Logf("ParseInt: %d\n", integer)
				}
			})
		}
	})

	t.Run("ParseByteString", func(t *testing.T) {
		tests := []DecoderTest{
			// valid
			{name: "valid bytestring", val: []byte("7:bencode"), fail: false},
			{name: "empty bytestring", val: []byte("0:"), fail: false},
			{name: "single character", val: []byte("1:a"), fail: false},
			{name: "classic spam", val: []byte("4:spam"), fail: false},
			{name: "with spaces", val: []byte("11:hello world"), fail: false},
			{name: "numeric content", val: []byte("10:0123456789"), fail: false},
			{name: "special characters", val: []byte("13:Hello, World!"), fail: false},
			{name: "punctuation only", val: []byte("3:..."), fail: false},
			{name: "url-like content", val: []byte("18:http://example.com"), fail: false},
			{name: "binary content", val: append([]byte("3:"), 0x00, 0x01, 0x02), fail: false},
			// invalid
			{name: "expect EOF", val: []byte("5:"), fail: true},
			{name: "length 1 no content", val: []byte("1:"), fail: true},
			{name: "truncated by one", val: []byte("5:hell"), fail: true},
			{name: "content shorter than length", val: []byte("10:abc"), fail: true},
			{name: "length much larger than content", val: []byte("100:short"), fail: true},
			{name: "negative length", val: []byte("-3:abc"), fail: true},
			{name: "missing length", val: []byte(":hello"), fail: true},
			{name: "no separator", val: []byte("5hello"), fail: true},
			{name: "letter prefix", val: []byte("abc"), fail: true},
			{name: "non-digit in length", val: []byte("3a:foo"), fail: true},
			{name: "leading zero in length", val: []byte("07:bencode"), fail: true},
			{name: "length overflow", val: []byte("99999999999999999999:x"), fail: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := Decoder{
					buf: make([]byte, 0),
					pos: 0,
				}

				d.buf = append(d.buf, tt.val...)

				bs, err := d.ParseByteString()
				if err != nil && !tt.fail {
					t.Errorf("ParseByteString %v: %v\n", tt.name, err)
				} else {
					t.Logf("ParseByteString %v: %v\n", tt.name, err)
				}

				if err == nil {
					fmt.Println("Parsed Byte String: ", string(bs))
				}
			})
		}
	})

	t.Run("ParseList", func(t *testing.T) {})
	t.Run("ParseDictionary", func(t *testing.T) {})
}
