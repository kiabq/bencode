package main

import "fmt"

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
func ParseThing(val []byte) (string, error) {
	switch {
	case val[0] >= '0' && val[0] <= '9':
		return "byte_string", nil
	case val[0] == 'i':
		return "integer", nil
	case val[0] == 'l':
		return "list", nil
	case val[0] == 'd':
		return "dictionary", nil
	default:
		return "", fmt.Errorf("Unknown Type: %v\n", val)
	}
}
