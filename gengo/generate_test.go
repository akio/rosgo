// Copyright 2018, Akio Ochiai All rights reserved
package main

import (
	"fmt"
	//	"math"
	"testing"
)

func TestGenerateMessage(t *testing.T) {
	const text string = `
# Comment
bool B = 1
int8 I8 =  -128
int16 I16 =  -32768
int32 I32 =  -2147483648
int64 I64 =  -9223372036854775808
uint8 U8 =  255
uint16 U16 =  65535
uint32 U32 =  4294967295
uint64 U64 =  18446744073709551615
string S = Lorem Ipsum # Comment is ignored

Header header
bool b
int8 i8
uint8 u8
int16 i16
uint16 u16
int32 i32
uint32 u32
int64 i64
uint64 u64
float32 f32
float64 f64
string s
time t
duration d
string[] sva
string[42] sfa
std_msgs/Empty e
std_msgs/Empty[] eva
std_msgs/Empty[42] efa
Bar x
Bar[] xva
Bar[42] xfa
`
	ctx, e := NewMsgContext()
	if e != nil {
		t.Errorf("Failed to create MsgContext.")
	}
	var spec *MsgSpec
	spec, e = LoadMsgFromString(ctx, text, "foo/Foo")
	if e != nil {
		t.Errorf("Failed to parse: %v", e)
	}

	msg, err := GenerateMessage(spec)
	if err != nil {
		t.Errorf("Failed to generate message: %v", err)
	}
	fmt.Printf(msg)
}
