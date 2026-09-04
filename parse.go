package main

import "fmt"

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

//	64, 32, 16, 8, 4, 2, 1
//
//	1, 0, 0, 0, 0, 1 = A (65 in ASCII)
//
// need a parser that can parse byte strings, integers, lists, and dictionaries

/*
* Byte Strings
for byte strings, we must guard against overflow on the length integer
we will want to determine if it's a valid byte string.
to determine that, you check if the byte string length is equal to the length value specified
*/
func ParseType(val []byte) (BencodeType, error) {
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
		return EncodedUnknown, fmt.Errorf("Unknown Type: %v\n", val)
	}
}

func ParseByteString() {}
func ParseInt()        {}
func ParseList()       {}
func ParseDictionary() {}
