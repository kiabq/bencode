package main

import (
	"errors"
	"fmt"
)

type BencodeType int

const (
	EncodedByteString BencodeType = iota
	EncodedInt
	EncodedList
	EncodedDictionary
	EncodedUnknown
)

func (b BencodeType) String() string {
	switch b {
	case EncodedByteString:
		return "byte_string"
	case EncodedInt:
		return "integer"
	case EncodedList:
		return "list"
	case EncodedDictionary:
		return "dictionary"
	default:
		return "unknown"
	}
}

// unsure how to structure the decoder struct...
// how to read strings into memeory?
// i think the idea is you read a file in chunks
// so you can stream to the decoder
type Decoder struct {
	buf []byte
	pos int
}

/*
64, 32, 16, 8, 4, 2, 1

1, 0, 0, 0, 0, 1 = A (65 in ASCII)

need a parser that can parse byte strings, integers, lists, and dictionaries

* Byte Strings
for byte strings, we must guard against overflow on the length integer
we will want to determine if it's a valid byte string.
to determine that, you check if the byte string length is equal to the length value specified
*/
func (d *Decoder) ParseType(val []byte) (BencodeType, error) {
	switch {
	case val[0] >= '0' && val[0] <= '9':
		return EncodedByteString, nil
	case val[0] == 'i':
		return EncodedInt, nil
	case val[0] == 'l':
		return EncodedList, nil
	case val[0] == 'd':
		return EncodedDictionary, nil
	default:
		return EncodedUnknown, fmt.Errorf("unknown Type: %v\n", val)
	}
}

func (d *Decoder) ParseByteString() ([]byte, error) {
	curr := d.buf[:d.pos]
	fmt.Printf("Current position (%v) and value (%v)\n", d.pos, curr)
	fmt.Printf("Value before cursor is %v\n", string(d.buf[:d.pos]))

	parsed, err := d.ParseType(curr)
	if err != nil {
		return nil, err
	}

	if parsed != EncodedByteString {
		return nil, errors.New("not valid byte string")
	}

	return curr, nil
}

func (d *Decoder) ParseInt() (int64, error) {
	return 0, nil
}

func (d *Decoder) ParseList() ([]byte, error) {
	return nil, nil
}

func (d *Decoder) ParseDictionary() ([]byte, error) {
	return nil, nil
}
