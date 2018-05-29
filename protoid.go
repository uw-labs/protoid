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

	r := &reader{buf: input}

	for !r.done() {
		val := r.decodeVarint()
		wiretype := val & 0x07
		k := int(val >> 3)
		switch wiretype {
		case 0: // varint value (int32, int64, uint32, uint64, sint32, sint64, bool, enum)
			v := r.decodeVarint()
			m[k] = v
		case 1: // 64 bit value (fixed64, sfixed64, double)
			v := r.readLeUint64()
			m[k] = v
		case 2: // length-delimited value (string, bytes, embedded messages, packed repeated fields)
			data := r.readLenDelimValue()

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
		case 3: // Start group (groups are deprecated)
			return nil, ErrNotImplemented
		case 4: // End group (groups are deprecated)
			return nil, ErrNotImplemented
		case 5: // 32-bit value (fixed32, sfixed32, float)
			val := r.readLeUint32()
			m[k] = val
		default:
			return nil, fmt.Errorf("unsupported wire type : %v", wiretype)
		}
	}
	if r.err != nil {
		return nil, r.err
	}
	return m, nil
}

type reader struct {
	buf []byte
	err error
}

func (r *reader) Err() error {
	return r.err
}

func (r *reader) done() bool {
	return len(r.buf) == 0 || r.err != nil
}

func (r *reader) readLeUint32() uint32 {
	if r.err != nil {
		return 0
	}

	if len(r.buf) < 4 {
		r.err = ErrUnexpectedEndOfInput
		return 0
	}
	v := binary.LittleEndian.Uint32(r.buf)
	r.buf = r.buf[4:]
	return v
}

func (r *reader) readLeUint64() uint64 {
	if r.err != nil {
		return 0
	}

	if len(r.buf) < 8 {
		r.err = ErrUnexpectedEndOfInput
		return 0
	}
	v := binary.LittleEndian.Uint64(r.buf)
	r.buf = r.buf[8:]
	return v
}

func (r *reader) readLenDelimValue() []byte {
	if r.err != nil {
		return nil // TODO: return empty slice?
	}

	l := r.decodeVarint()

	if uint64(len(r.buf)) < l {
		r.err = ErrUnexpectedEndOfInput
		return nil // TODO: return empty slice?
	}
	data := r.buf[0:l]
	r.buf = r.buf[l:]
	return data
}

func (r *reader) decodeVarint() uint64 {
	if r.err != nil {
		return 0
	}

	var l int
	var val uint64
	for shift := uint(0); shift < 64; shift += 7 {
		if l >= len(r.buf) {
			r.err = ErrUnexpectedEndOfInput
			return 0
		}
		b := uint64(r.buf[l])
		l++
		val |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			r.buf = r.buf[l:]
			return val
		}
	}

	r.err = ErrNumberTooLarge
	return 0

}

func copyBytes(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
