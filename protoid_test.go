package protoid

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)

	ss := &SingleString{TheString: "string123"}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: "string123"}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestBytes(t *testing.T) {

	assert := assert.New(t)

	ss := &SingleBytes{TheBytes: []byte{255, 0, 77, 66, 55}}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: []byte{255, 0, 77, 66, 55}}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestInt32(t *testing.T) {
	assert := assert.New(t)

	ss := &SingleInt32{TheInt32: 123456}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: uint64(123456)}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)

	ss := &SingleBool{TheBool: true}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: uint64(1)}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestTwoStrings(t *testing.T) {
	assert := assert.New(t)

	ss := &TwoStrings{String_1: "string1", String_2: "string2"}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: "string1", 2: "string2"}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestRepeatedString(t *testing.T) {
	assert := assert.New(t)

	ss := &RepeatedString{MyString: []string{"A", "B", "C"}}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: []interface{}{"A", "B", "C"}}

	t.Logf("%x\n", ser)
	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestRepeatedInt32(t *testing.T) {
	t.Skip("this doesn't work yet. We need more intelligence")

	assert := assert.New(t)

	ss := &RepeatedInt32{MyInt32S: []int32{1, 2, 3}}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: []interface{}{1, 2, 3}}

	t.Logf("%x\n", ser)
	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestSingleEnum(t *testing.T) {
	assert := assert.New(t)

	ss := &SingleEnum{TheEnum: TestEnum_VAL_1}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: uint64(1)}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestSingleFixed64(t *testing.T) {
	assert := assert.New(t)

	ss := &SingleFixed64{TheFixed64: 12345678}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: uint64(12345678)}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestSingleFixed32(t *testing.T) {
	assert := assert.New(t)

	ss := &SingleFixed32{TheFixed32: 12345678}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: uint32(12345678)}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestSingleEmbedded(t *testing.T) {

	assert := assert.New(t)

	ss := &SingleEmbedded{MySingleString: &SingleString{TheString: "123"}}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: map[int]interface{}{1: "123"}}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}

func TestRepeatedEmbedded(t *testing.T) {

	assert := assert.New(t)

	ss := &RepeatedEmbedded{MySingleStrings: []*SingleString{
		&SingleString{TheString: "123"}, &SingleString{TheString: "456"},
	}}
	ser, err := proto.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[int]interface{}{1: []interface{}{
		map[int]interface{}{1: "123"},
		map[int]interface{}{1: "456"},
	}}

	actual, err := Decode(ser)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected, actual)
}
