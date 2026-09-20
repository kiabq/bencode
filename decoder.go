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

// TODO: PASS IO READER INTO THIS METHOD
func (d *Decoder) Decode() error {
	b := d.buf[d.pos]

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
	for i, b := range d.buf[d.pos:] {
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
			integerSlice := d.buf[d.pos+1 : d.pos+i]
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

			d.pos = d.pos + i + 1

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

	sep := false
	byteString := make([]byte, 0)

	for i, b := range d.buf[d.pos:] {

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

					sentinel = d.buf[d.pos : d.pos+i]
					v, err := strconv.ParseInt(string(sentinel), 10, 64)
					if err != nil {
						if errors.Is(err, strconv.ErrRange) {
							return nil, errors.New("length had int overflow")
						}

						// error something went wrong parsing length (EOF error)
					}

					length = int(v)

					d.pos = d.pos + len(sentinel) + 1

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
					d.pos = d.pos + len(byteString)
					return byteString, nil
				}
			}
		}
	}

	// EOF
	if len(byteString) != length {
		// TODO MAKE ERRORS IN ERROR.GO
		return byteString, errors.New("EOF")
	}

	return nil, nil
}

func (d *Decoder) ParseList() ([]interface{}, error) {
	var list []interface{}
	var term bool

	if len(d.buf) == 0 {
		return nil, errors.New("fucked up input my brother")
	}

	for d.pos < len(d.buf) {
		if d.buf[d.pos] == 'l' {
			d.pos++

			if d.pos > 1 {
				nestedList, err := d.ParseList()

				list = append(list, nestedList)

				if err != nil {
					return nil, err
				}
			}

			continue
		}

		// last e in list, if no error by now, valid terminator
		if d.buf[d.pos] == 'e' {
			d.pos++
			term = true
			break
		}

		if d.buf[d.pos] > '0' && d.buf[d.pos] < '9' {
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
		return nil, errors.New("no terminator")
	}

	return list, nil
}

func (d *Decoder) ParseDictionary() ([]byte, error) {
	d.pos++

	return nil, nil
}
