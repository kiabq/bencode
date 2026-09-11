package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
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

func (d *Decoder) GetPieceLength(piece []byte) (int64, error) {
	// len is type int, not sure how to determine bit size on any given arch.
	length := len(piece)

	if length > math.MaxInt {
		return -1, errors.New("piece length exceeds int64 max value")
	}

	return 0, nil
}

func (d *Decoder) ParseByteString() ([]byte, error) {
	curr := d.buf[d.pos:]
	var k, v []byte
	sep := false

	for _, b := range curr {
		d.pos++
		// could be found before string is finished computing
		// should handle that as an error, not now though
		// assuming all strings and lengths are equal for now
		if string(b) == ":" {
			sep = true
			continue
		}

		if sep {
			v = append(v, b)
		} else {
			k = append(k, b)
		}
	}

	length, err := strconv.Atoi(string(k))
	if err != nil {
		return nil, err
	}

	fmt.Println(string(v), length)

	if length != len(v) {
		return nil, errors.New("ParseByteString failed to decode")
	}

	return v, nil
}

func (d *Decoder) ParseInt() (int64, error) {
	//curr := d.buf[d.pos:]

	return 0, nil
}

func (d *Decoder) ParseList() ([]byte, error) {
	curr := d.buf[d.pos:]

	for _, b := range curr {
		if b == 'l' {
			d.pos++
			continue
		} else if b == 'e' {
			continue
		} else {
			v, err := d.ParseByteString()
			if err != nil {
				return nil, err
			}
			fmt.Println("ParseList() ", string(v))
		}
	}

	return nil, nil
}

func (d *Decoder) ParseDictionary() ([]byte, error) {
	//curr := d.buf[d.pos:]

	return nil, nil
}
