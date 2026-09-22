package main

// POSSIBLE ERRORS
// Null root value          - :<val>, missing <key>
// Non-singular root item   - 4:spami42e, i42e is another element, but string was parsed.
// Invalid type encountered - not i, l, d, or 0-9
// Missing 'e' terminator   - on i, l, d
// Integer errors
//	Contains non-digit characters
//  Has a leading zero
//  Is negative zero
//  Int Underflow
//  Int Overflow
// Byte string errors
//	Negative length
// 	Length not followed by ':'
//  Unexpected EOF before completing string
//  Length specified in units of codepoints (characters) rather than bytes.
// Dictionary errors
// 	Key is not string
//  Duplicate keys
//  Keys not sorted
//  Keys incorrectly sorted by codepoint in a particular character encoding, rather than lexicographically sorted by ordinal
// 	Missing value for a key
