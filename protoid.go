package protoid

import (
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

type valueApplier interface {
	mapType0(propnum int, value uint64) error
	mapType1(propnum int, value uint64) error
	mapType2(propnum int, value []byte) error
	mapType5(propnum int, value uint32) error
}

type genericMapValueApplier struct {
	m map[int]interface{}
}

func (va *genericMapValueApplier) mapType0(propnum int, value uint64) error {
	va.m[propnum] = value
	return nil
}

func (va *genericMapValueApplier) mapType1(propnum int, value uint64) error {
	va.m[propnum] = value
	return nil
}

func (va *genericMapValueApplier) mapType2(propnum int, data []byte) error {

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

	if va.m[propnum] != nil {
		// we already have a value here, so this must be a repeated value.
		slice, ok := va.m[propnum].([]interface{})
		if ok {
			// already a slice, simply append.
			slice = append(slice, value)
		} else {
			// single value currently, change to a slice.
			slice = append(slice, va.m[propnum])
			slice = append(slice, value)
		}
		va.m[propnum] = slice
	} else {
		va.m[propnum] = value
	}
	return nil
}

func (va *genericMapValueApplier) mapType5(propnum int, value uint32) error {
	va.m[propnum] = value
	return nil
}

// Decode decodes an arbitrary protocol buffers message into a map of field number to field value. It makes a best-effort attempt to use the most appropriate type for the values.  Embedded structs, strings, integers and more are often decoded correctly.  However due to the nature of protocol buffers, it is not always possible to do this perfectly.
func Decode(input []byte) (map[int]interface{}, error) {
	m := make(map[int]interface{})

	va := &genericMapValueApplier{m}

	if err := decode(input, va); err != nil {
		return nil, err
	}

	return m, nil
}

func decode(input []byte, va valueApplier) error {

	r := &reader{buf: input}

	for !r.done() {
		val := r.decodeVarint()
		wiretype := val & 0x07
		k := int(val >> 3)
		switch wiretype {
		case 0: // varint value (int32, int64, uint32, uint64, sint32, sint64, bool, enum)
			v := r.decodeVarint()
			if err := va.mapType0(k, v); err != nil {
				return err
			}
		case 1: // 64 bit value (fixed64, sfixed64, double)
			v := r.readLeUint64()
			if err := va.mapType1(k, v); err != nil {
				return err
			}
		case 2: // length-delimited value (string, bytes, embedded messages, packed repeated fields)
			v := r.readLenDelimValue()
			if err := va.mapType2(k, v); err != nil {
				return err
			}
		case 3: // Start group (groups are deprecated)
			return ErrNotImplemented
		case 4: // End group (groups are deprecated)
			return ErrNotImplemented
		case 5: // 32-bit value (fixed32, sfixed32, float)
			v := r.readLeUint32()
			if err := va.mapType5(k, v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported wire type : %v", wiretype)
		}
	}
	if r.err != nil {
		return r.err
	}
	return nil
}

func copyBytes(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
