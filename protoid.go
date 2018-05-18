package protoid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf8"
)

var (
	// ErrUnexpectedEndOfInput indicates that the input data is shorter than expected.
	ErrUnexpectedEndOfInput = errors.New("unexpected end of input")
	// ErrNotImplemented indicates that a message contained somethings that protoid cannot currently decode.
	ErrNotImplemented = errors.New("groups not implemented")
	// ErrNumberTooLarge inducates that a number found in a protocol buffers message cannot be represented as a 64 bit value.
	ErrNumberTooLarge = errors.New("number too large for 64 bit value")
)

// Decode decodes an arbitrary protocol buffers message into a map of field number to field value. It makes a best-effort attempt to use the most appropriate type for the values.  Embedded structs, strings, integers and more are often decoded correctly.  However due to the nature of protocol buffers, it is not always possible to do this perfectly.
func Decode(input []byte) (map[int]interface{}, error) {
	m := make(map[int]interface{})
	for len(input) > 0 {
		val, i, err := decodeVarint(input)
		if err != nil {
			return nil, err
		}
		input = input[i:]
		wiretype := val & 0x07
		k := int(val >> 3)
		switch wiretype {
		case 0: // varint value (int32, int64, uint32, uint64, sint32, sint64, bool, enum)
			v, i, err := decodeVarint(input)
			if err != nil {
				return nil, err
			}
			input = input[i:]
			m[k] = v
		case 1: // 64 bit value (fixed64, sfixed64, double)
			if len(input) < 8 {
				return nil, ErrUnexpectedEndOfInput
			}
			v := binary.LittleEndian.Uint64(input)
			input = input[8:]
			m[k] = v
		case 2: // length-delimited value (string, bytes, embedded messages, packed repeated fields)
			l, i, err := decodeVarint(input)
			if err != nil {
				return nil, err
			}
			input = input[i:]

			if uint64(len(input)) < l {
				return nil, ErrUnexpectedEndOfInput
			}
			data := input[0:l]

			var value interface{}
			// try to guess the type of data
			// first try to decode as embedded value
			emb, err := Decode(data)
			if err == nil {
				value = emb
			} else {
				if utf8.Valid(data) {
					// string
					value = string(data)
				} else {
					// assume bytes
					value = copyBytes(data)
				}
			}

			if m[k] != nil {
				// we already have a value here, so this must be a repeated value.
				slice, ok := m[k].([]interface{})
				if ok {
					// already a slice, simply append.
					slice = append(slice, value)
				} else {
					// single value currently, change to a slice.
					slice = append(slice, m[k])
					slice = append(slice, value)
				}
				m[k] = slice
			} else {
				m[k] = value
			}
			if uint64(len(input)) < l {
				return nil, ErrUnexpectedEndOfInput
			}
			input = input[l:]
		case 3: // Start group (groups are deprecated)
			return nil, ErrNotImplemented
		case 4: // End group (groups are deprecated)
			return nil, ErrNotImplemented
		case 5: // 32-bit value (fixed32, sfixed32, float)
			if len(input) < 4 {
				return nil, ErrUnexpectedEndOfInput
			}
			val := binary.LittleEndian.Uint32(input)
			input = input[4:]
			m[k] = val
		default:
			return nil, fmt.Errorf("unsupported wire type : %v\n", wiretype)
		}
	}
	return m, nil
}

func decodeVarint(buf []byte) (val uint64, l int, err error) {
	for shift := uint(0); shift < 64; shift += 7 {
		if l >= len(buf) {
			return 0, 0, ErrUnexpectedEndOfInput
		}
		b := uint64(buf[l])
		l++
		val |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return val, l, nil
		}
	}

	return 0, 0, ErrNumberTooLarge
}

func copyBytes(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
