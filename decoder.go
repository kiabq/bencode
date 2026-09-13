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
}

func (d *Decoder) Decode() error {
	b := d.buf[d.pos]

	switch b {
	case 'i':
		_, err := d.ParseInt()
		if err != nil {
			return err
		}
	case 'l':
		return nil
	case 'd':
		return nil
	default:
		if b >= '0' && b <= '9' {
			return nil
		}
	}

	return nil
}

func (d *Decoder) ParseInt() (int64, error) {
	startZero := false
	startNegative := false
	for i, b := range d.buf[d.pos:] {
		if i == 0 { // byte is 'i', skip loop
			continue
		}

		if i == 1 {
			if b == '0' { // check for zero, used to determine leading zero later
				startZero = true
			}

			if b == '-' { // check if int negative
				startNegative = true
				continue
			}
		} else if b >= '0' && b <= '9' && !startZero { // then, check if is number
			if b == '0' && startNegative {
				return -1, errors.New("negative zero")
			}

			continue
		} else if b == 'e' { // if not number, then it is end.
			integerSlice := d.buf[d.pos+1 : d.pos+i]
			integer, err := strconv.ParseInt(string(integerSlice), 10, 64)
			if err != nil {
				// err int overflow
				if errors.Is(err, strconv.ErrRange) {
					// maybe special error?
				}

				return -1, err // parse int failed
			} else {
				d.pos += i + 1
				return integer, nil
			}

		} else {
			if startZero {
				return -1, errors.New("leading zero error")
			}
		}
	}

	return 0, nil
}

func (d *Decoder) ParseByteString() ([]byte, error) {
	//var length string
	for i, b := range d.buf[d.pos:] {
		if i == 0 {
			if b == '-' { // invalid, byte string can't be negative
				// throw error for negative byte string
			}

			if b == ':' {
				// no length, error
			}
		}
	}

	return nil, nil
}

func (d *Decoder) ParseList() ([]byte, error) {
	d.pos++

	return nil, nil
}

func (d *Decoder) ParseDictionary() ([]byte, error) {
	d.pos++

	return nil, nil
}
