package main

import (
	"errors"
	"fmt"
	"strconv"
)

/* example of structured bencode

d:                        {meaning: 42, disturb: [bencode, -20]}
	7:meaning:
		42
	7:disturb:
		l7:               [bencode, -20]
			bencode
			i-20e
		e
e
*/

type Decoder struct {
	buf  []byte
	pos  int
	step int
	Ret  interface{} `json:"ret"` // how the fuck do i build this?
}

func (d *Decoder) curr() byte {
	return d.buf[d.pos]
}

func (d *Decoder) read() []byte {
	return d.buf[d.pos:]
}

func (d *Decoder) slice(start, end int) ([]byte, error) {
	lo, hi := d.pos+start, d.pos+end
	if lo < 0 || hi > len(d.buf) || lo > hi {
		return nil, errors.New("slice out of bounds")
	}
	return d.buf[lo:hi], nil
}

func (d *Decoder) advance(offset int) {
	d.step++
	d.pos = d.pos + offset
}

// TODO: PASS IO READER INTO THIS METHOD
func (d *Decoder) Decode() error {
	b := d.curr()

	switch b {
	case 'i':
		_, err := d.ParseInt()
		if err != nil {
			return err
		}
	case 'l':
		_, err := d.ParseList()
		if err != nil {
			return err
		}
	case 'd':
		_, err := d.ParseDictionary()
		if err != nil {
			return err
		}
	default:
		_, err := d.ParseByteString() // figure out wtf to do with the value of this func
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) ParseInt() (int64, error) {
	startZero := false
	startNegative := false
	for i, b := range d.read() {
		if i == 0 { // byte is 'i', skip loop
			continue
		}

		if i == 1 {
			switch b {
			case '0':
				startZero = true
			case '-':
				startNegative = true
			case '+':
				return 0, errors.New("incorrect prefix")
			case 'e':
				return 0, errors.New("no integer")
			}

			continue
		}

		if b >= '0' && b <= '9' { // then, check if is number
			if startZero {
				return 0, errors.New("leading zero error")
			}

			if i == 2 && b == '0' && startNegative {
				return 0, errors.New("negative zero")
			}

			continue
		}

		if b == 'e' { // if not number, then it is end.
			integerSlice, err := d.slice(1, i)
			if err != nil {
				return 0, err
			}

			integer, err := strconv.ParseInt(string(integerSlice), 10, 64)

			// err int underflow
			if startNegative && integer > 0 {
				return 0, errors.New("integer underflow")
			}

			if err != nil {
				// err int overflow
				if errors.Is(err, strconv.ErrRange) {
					// maybe special error?
					return 0, errors.New("integer overflow")
				}

				return 0, err // parse int failed
			}

			d.advance(i + 1) // account for length of current []byte position + next value 'e'

			return integer, nil
		} else {
			return 0, errors.New("wrong terminator")
		}
	}

	return 0, errors.New("unexpected end of input")
}

func (d *Decoder) ParseByteString() ([]byte, error) {
	index := d.pos
	byteString := make([]byte, 0)
	byteSize := int64(0)

	// scan to within bound of buffer
	// a recursive descent parser would work if the grammar
	// could be constructed as tree, but grammar can't be dynamic
	// the best we can do is peek ahead

	// use parseint when we find value
	for d.pos < len(d.buf) {
		if d.curr() == ':' && byteSize == 0 {
			sentinel := string(d.buf[index:d.pos])

			// get bytestring length
			length, err := strconv.ParseInt(sentinel, 10, 64)
			if err != nil {
				if errors.Is(err, strconv.ErrRange) {
					return nil, errors.New("length had int overflow")
				} else {
					// not quite sure what to do with other errors atm
					return nil, err
				}
			}

			// leading zero error
			if sentinel[0] == '0' && len(sentinel) >= 2 {
				return nil, errors.New("leading zero error")
			}

			d.advance(1)

			if length == 0 {
				return []byte(""), nil
			} else if length > 0 {
				byteSize = length
			} else {
				// catch negative length
				return nil, errors.New("negative length")
			}
			continue
		}

		if byteSize == 0 {
			d.advance(1)
			continue
		}

		byteString = append(byteString, d.curr())

		if int64(len(byteString)) == byteSize {
			d.advance(1)
			return byteString, nil
		}

		d.advance(1)
	}

	return nil, errors.New("unexpected end of input")
}

func (d *Decoder) ParseList() ([]interface{}, error) {
	var list []interface{}
	var term bool

	d.advance(1) // eat current byte because current byte sits on 'l'. prevent infinite loops

	// TODO ACCOUNT FOR DICTIONARIES IN LIST
	for d.pos < len(d.buf) {
		if d.curr() == 'e' { // compliments d.advance(1) at top of func. if 'e' is found, it belongs to the list
			d.advance(1)
			term = true
			break
		}

		if d.curr() == 'd' {
			continue
		} else if d.curr() == 'l' {
			nestedList, err := d.ParseList()
			if err != nil {
				return nil, err
			}
			list = append(list, nestedList)
		} else if d.curr() >= '0' && d.curr() <= '9' {
			bs, err := d.ParseByteString()
			if err != nil {
				return nil, err
			}
			list = append(list, bs)
		} else {
			integer, err := d.ParseInt()
			if err != nil {
				return nil, err
			}
			list = append(list, integer)
		}
	}

	if !term {
		return nil, errors.New("ParseList - no terminator")
	}

	return list, nil
}

func (d *Decoder) ParseDictionary() ([]byte, error) {
	d.pos++

	return nil, nil
}
