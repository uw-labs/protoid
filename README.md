Protoid
=======

[![go-doc](https://godoc.org/github.com/uw-labs/protoid?status.svg)](https://godoc.org/github.com/uw-labs/protoid)

Protoid is a lightweight simple protocol buffers decoder. It is intended to make as much sense as possible of a protocol buffers message without any definitions being available.

Because field names are unknown, we only have the field numbers in the output. Where possible values will be correctly represented as their proper type : strings, bytes, embedded structs, integers etc.  []interface{} is used for all repeated values.

All current proto3 data is supported. proto2 support is incomplete.

Limitations
-----------
Due to the design choices of protobuf, it is understandably impossible to always correctly know the type of the values.  protoid can only make a best effort guess by inspecting the data.
