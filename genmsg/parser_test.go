// Copyright 2018, Akio Ochiai All rights reserved
package genmsg

import (
	"fmt"
	//	"math"
	"reflect"
	"testing"
)

func TestConvertConstantValue(t *testing.T) {
	var tests = []struct {
		fieldType    string
		valueLiteral string
		expected     interface{}
		expectError  bool
	}{
		{"bool", "0", false, false},
		{"bool", "1", true, false},
		{"bool", "2", true, false},
		{"bool", "-2", true, true},
		{"bool", "True", true, false},
		{"bool", "False", false, false},
		{"bool", "None", false, false},
		{"float32", "2.72", float32(2.72), false},
		{"float64", "-3.14", float64(-3.14), false},
		{"int8", "-129", 0, true},
		{"int8", "-128", int8(-128), false},
		{"int8", "127", int8(127), false},
		{"int8", "128", 0, true},
		{"int16", "-32769", 0, true},
		{"int16", "-32768", int16(-32768), false},
		{"int16", "32767", int16(32767), false},
		{"int16", "32768", 0, true},
		{"int32", "-2147483649", 0, true},
		{"int32", "-2147483648", int32(-2147483648), false},
		{"int32", "2147483647", int32(2147483647), false},
		{"int32", "2147483648", 0, true},
		{"int64", "-9223372036854775809", 0, true},
		{"int64", "-9223372036854775808", int64(-9223372036854775808), false},
		{"int64", "9223372036854775807", int64(9223372036854775807), false},
		{"int64", "9223372036854775808", 0, true},
		{"uint8", "-1", 0, true},
		{"uint8", "0", uint8(0), false},
		{"uint8", "255", uint8(255), false},
		{"uint8", "256", 0, true},
		{"uint16", "-1", 0, true},
		{"uint16", "0", uint16(0), false},
		{"uint16", "65535", uint16(65535), false},
		{"uint16", "65536", 0, true},
		{"uint32", "-1", 0, true},
		{"uint32", "0", uint32(0), false},
		{"uint32", "4294967295", uint32(4294967295), false},
		{"uint32", "4294967296", 0, true},
		{"uint64", "-1", 0, true},
		{"uint64", "0", uint64(0), false},
		{"uint64", "18446744073709551615", uint64(18446744073709551615), false},
		{"uint64", "18446744073709551616", 0, true},
		{"string", "Lorem Ipsum", "Lorem Ipsum", false},
	}

	for _, test := range tests {
		result, e := convertConstantValue(test.fieldType, test.valueLiteral)
		if test.expectError {
			if e == nil {
				t.Errorf("INPUT(%s : %s) | should fail but succeeded", test.valueLiteral, test.fieldType)
			}
		} else {
			if e != nil {
				t.Errorf("INPUT(%s : %s) | %s", test.valueLiteral, test.fieldType, e.Error())

			} else if result != test.expected {
				format := "INPUT(%s : %s) | Expected: [%v: %v], Actual: [%v : %v]"
				t.Errorf(format, test.valueLiteral, test.fieldType, test.expected, reflect.TypeOf(test.expected), result, reflect.TypeOf(result))
			}
		}
	}
}

func TestParseMessage(t *testing.T) {
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
	fmt.Println("---")
	fmt.Println(spec.String())

}

func TestParseService(t *testing.T) {

}
