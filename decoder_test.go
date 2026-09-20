package main

import (
	"reflect"
	"testing"
)

type DecoderTest struct {
	name string
	val  []byte
	pos  int         // starting position within buf
	fail bool        // true if a parse error is expected
	want interface{} // expected decoded value on success (nil = skip value check)
	end  int         // expected d.pos after a successful parse (0 = skip)
}

func TestDecoder(t *testing.T) {
	t.Run("ParseInt", func(t *testing.T) {
		tests := []DecoderTest{
			// valid
			{name: "positive integer", val: []byte("i105e"), want: int64(105), end: 5},
			{name: "zero integer", val: []byte("i0e"), want: int64(0), end: 3},
			{name: "single digit", val: []byte("i7e"), want: int64(7), end: 3},
			{name: "in a list", val: []byte("i42e"), pos: 3, want: int64(42), end: 7},
			{name: "negative integer", val: []byte("i-42e"), want: int64(-42), end: 5},
			{name: "negative single digit", val: []byte("i-1e"), want: int64(-1), end: 4},
			{name: "max int64", val: []byte("i9223372036854775807e"), want: int64(9223372036854775807), end: 21},
			{name: "min int64", val: []byte("i-9223372036854775808e"), want: int64(-9223372036854775808), end: 22},
			{name: "stops at terminator", val: []byte("i42ejunk"), want: int64(42), end: 4},
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
			{name: "spaces inside", val: []byte("i 42 e"), fail: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := Decoder{
					buf: make([]byte, tt.pos),
					pos: tt.pos,
				}
				d.buf = append(d.buf, tt.val...)

				got, err := d.ParseInt()
				if tt.fail {
					if err == nil {
						t.Fatalf("expected error, got value %d (pos %d)", got, d.pos)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.want != nil && got != tt.want.(int64) {
					t.Errorf("value = %d, want %d", got, tt.want)
				}
				if tt.end != 0 && d.pos != tt.end {
					t.Errorf("end pos = %d, want %d", d.pos, tt.end)
				}
			})
		}
	})

	t.Run("ParseByteString", func(t *testing.T) {
		tests := []DecoderTest{
			// valid
			{name: "valid bytestring", val: []byte("7:bencode"), want: "bencode", end: 9},
			{name: "empty bytestring", val: []byte("0:"), want: "", end: 2},
			{name: "single character", val: []byte("1:a"), want: "a", end: 3},
			{name: "classic spam", val: []byte("4:spam"), want: "spam", end: 6},
			{name: "with spaces", val: []byte("11:hello world"), want: "hello world", end: 14},
			{name: "numeric content", val: []byte("10:0123456789"), want: "0123456789", end: 13},
			{name: "special characters", val: []byte("13:Hello, World!"), want: "Hello, World!", end: 16},
			{name: "punctuation only", val: []byte("3:..."), want: "...", end: 5},
			{name: "url-like content", val: []byte("18:http://example.com"), want: "http://example.com", end: 21},
			{name: "binary content", val: append([]byte("3:"), 0x00, 0x01, 0x02), want: string([]byte{0x00, 0x01, 0x02}), end: 5},
			{name: "accounts for spaces", val: []byte("12:hello  world"), want: "hello  world", end: 15}, // 2 spaces
			{name: "stops at declared length", val: []byte("4:spamiextra"), want: "spam", end: 6},
			{name: "digit content matches length", val: []byte("1:5"), want: "5", end: 3},
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
					pos: tt.pos,
				}
				d.buf = append(d.buf, tt.val...)

				got, err := d.ParseByteString()
				if tt.fail {
					if err == nil {
						t.Fatalf("expected error, got %q (pos %d)", got, d.pos)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.want != nil && string(got) != tt.want.(string) {
					t.Errorf("value = %q, want %q", got, tt.want)
				}
				if tt.end != 0 && d.pos != tt.end {
					t.Errorf("end pos = %d, want %d", d.pos, tt.end)
				}
			})
		}
	})

	t.Run("ParseList", func(t *testing.T) {
		tests := []DecoderTest{
			// valid
			{name: "empty list", val: []byte("le"), end: 2},
			{name: "single integer", val: []byte("li42ee"), want: []interface{}{int64(42)}, end: 6},
			{name: "single string", val: []byte("l4:spame"), want: []interface{}{[]byte("spam")}, end: 8},
			{name: "multiple integers", val: []byte("li1ei2ei3ee"), want: []interface{}{int64(1), int64(2), int64(3)}, end: 11},
			{name: "multiple strings", val: []byte("l4:spam3:fooe"), want: []interface{}{[]byte("spam"), []byte("foo")}, end: 13},
			{name: "mixed types", val: []byte("l4:spami42ee"), want: []interface{}{[]byte("spam"), int64(42)}, end: 12},
			{name: "negative integer", val: []byte("li-5ee"), want: []interface{}{int64(-5)}, end: 6},
			{name: "element after nested list", val: []byte("lli42eei7ee"), want: []interface{}{[]interface{}{int64(42)}, int64(7)}, end: 11},
			{name: "nested list", val: []byte("lli42eee"), want: []interface{}{[]interface{}{int64(42)}}, end: 8},
			{name: "nested empty list", val: []byte("llee"), end: 4},
			{name: "string with spaces", val: []byte("l11:hello worlde"), want: []interface{}{[]byte("hello world")}, end: 16},
			// invalid
			{name: "missing terminator", val: []byte("li42e"), fail: true},
			{name: "unclosed nested list", val: []byte("lli42ee"), fail: true},
			{name: "invalid integer inside", val: []byte("li042ee"), fail: true},
			{name: "invalid string inside", val: []byte("l5:hie"), fail: true},
			{name: "empty input", val: []byte(""), fail: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d := Decoder{
					buf: make([]byte, 0),
					pos: tt.pos,
					Ret: nil,
				}
				d.buf = append(d.buf, tt.val...)

				got, err := d.ParseList()
				if tt.fail {
					if err == nil {
						t.Fatalf("expected error, got %#v (pos %d)", got, d.pos)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("value = %#v, want %#v", got, tt.want)
				}
				if tt.end != 0 && d.pos != tt.end {
					t.Errorf("end pos = %d, want %d", d.pos, tt.end)
				}
			})
		}
	})

	t.Run("ParseDictionary", func(t *testing.T) {})
}
