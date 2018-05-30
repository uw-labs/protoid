package protoid

import "encoding/binary"

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
