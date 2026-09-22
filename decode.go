package main

import (
	"errors"
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
	buf []byte
	pos int
	Ret interface{} `json:"ret"` // how the fuck do i build this?
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

func (d *Decoder) expand(offset int) {
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

			d.expand(i + 1) // account for length of current []byte position + next value 'e'

			return integer, nil
		} else {
			return 0, errors.New("wrong terminator")
		}
	}

	return 0, errors.New("unexpected end of input")
}

func (d *Decoder) ParseByteString() ([]byte, error) {
	var length int
	var sentinel []byte
	var leadingZero bool
	var err error

	sep := false
	byteString := make([]byte, 0)

	// WIP refactor
	for d.pos < len(d.buf) {

		// check if - ;  < 0 or > 9
	}

	for i, b := range d.read() {
		if i == 0 {
			if b == '-' { // invalid, byte string can't be negative
				// throw error for negative byte string
				return nil, errors.New("length can't be negative")
			}

			if b == ':' || b < '0' || b > '9' {
				// no length, error
				return nil, errors.New("no length")
			}

			if b == '0' {
				leadingZero = true
			}
		} else {
			if !sep {
				if b == ':' {
					sep = true

					sentinel, err = d.slice(0, i)
					if err != nil {
						return nil, err
					}

					v, err := strconv.ParseInt(string(sentinel), 10, 64)
					if err != nil {
						if errors.Is(err, strconv.ErrRange) {
							return nil, errors.New("length had int overflow")
						} else {
							// not quite sure what to do with other errors atm
							return nil, err
						}
					}

					length = int(v)
					d.expand(len(sentinel) + 1) // account for length of int slice + the semicolon
					continue
				} else {
					if leadingZero {
						return nil, errors.New("length can't have leading zero")
					}

					if b < '0' || b > '9' {
						// error, length not followed by separator
						return nil, errors.New("missing separator")
					}
				}
			} else {
				byteString = append(byteString, b)

				if len(byteString) == length {
					d.expand(len(byteString)) // advance cursor by length of string
					return byteString, nil
				}
			}
		}
	}

	if length == 0 && leadingZero && sep {
		return []byte(""), nil
	}

	// EOF
	return byteString, errors.New("EOF")
}

func (d *Decoder) ParseList() ([]interface{}, error) {
	var list []interface{}
	var term bool

	d.expand(1) // eat current byte because current byte sits on 'l'. prevent infinite loops

	// TODO ACCOUNT FOR DICTIONARIES IN LIST
	for d.pos < len(d.buf) {
		if d.curr() == 'e' { // compliments d.expand(1) at top of func. if 'e' is found, it belongs to the list
			d.expand(1)
			term = true
			break
		}

		if d.curr() == 'l' {
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
