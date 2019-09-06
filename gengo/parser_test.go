// Copyright 2018, Akio Ochiai All rights reserved
package main

import (
	// "fmt"
	//	"math"
	"os"
	"reflect"
	"strings"
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

	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")
	ctx, e := NewMsgContext(strings.Split(rosPkgPath, ":"))
	if e != nil {
		t.Errorf("Failed to create MsgContext.")
	}
	// var spec *MsgSpec
	_, e = ctx.LoadMsgFromString(text, "foo/Foo")
	if e != nil {
		t.Errorf("Failed to parse: %v", e)
	}
	// fmt.Println("---")
	// fmt.Println(spec.String())
}

func TestMD5_std_msgs(t *testing.T) {
	var std_msgs = map[string]string{
		"std_msgs/Bool":                "8b94c1b53db61fb6aed406028ad6332a",
		"std_msgs/Byte":                "ad736a2e8818154c487bb80fe42ce43b",
		"std_msgs/ByteMultiArray":      "70ea476cbcfd65ac2f68f3cda1e891fe",
		"std_msgs/Char":                "1bf77f25acecdedba0e224b162199717",
		"std_msgs/ColorRGBA":           "a29a96539573343b1310c73607334b00",
		"std_msgs/Duration":            "3e286caf4241d664e55f3ad380e2ae46",
		"std_msgs/Empty":               "d41d8cd98f00b204e9800998ecf8427e",
		"std_msgs/Float32":             "73fcbf46b49191e672908e50842a83d4",
		"std_msgs/Float32MultiArray":   "6a40e0ffa6a17a503ac3f8616991b1f6",
		"std_msgs/Float64":             "fdb28210bfa9d7c91146260178d9a584",
		"std_msgs/Float64MultiArray":   "4b7d974086d4060e7db4613a7e6c3ba4",
		"std_msgs/Header":              "2176decaecbce78abc3b96ef049fabed",
		"std_msgs/Int16":               "8524586e34fbd7cb1c08c5f5f1ca0e57",
		"std_msgs/Int16MultiArray":     "d9338d7f523fcb692fae9d0a0e9f067c",
		"std_msgs/Int32":               "da5909fbe378aeaf85e547e830cc1bb7",
		"std_msgs/Int32MultiArray":     "1d99f79f8b325b44fee908053e9c945b",
		"std_msgs/Int64":               "34add168574510e6e17f5d23ecc077ef",
		"std_msgs/Int64MultiArray":     "54865aa6c65be0448113a2afc6a49270",
		"std_msgs/Int8":                "27ffa0c9c4b8fb8492252bcad9e5c57b",
		"std_msgs/Int8MultiArray":      "d7c1af35a1b4781bbe79e03dd94b7c13",
		"std_msgs/MultiArrayDimension": "4cd0c83a8683deae40ecdac60e53bfa8",
		"std_msgs/MultiArrayLayout":    "0fed2a11c13e11c5571b4e2a995a91a3",
		"std_msgs/String":              "992ce8a1687cec8c8bd883ec73ca41d1",
		"std_msgs/Time":                "cd7166c74c552c311fbcc2fe5a7bc289",
		"std_msgs/UInt16":              "1df79edf208b629fe6b81923a544552d",
		"std_msgs/UInt16MultiArray":    "52f264f1c973c4b73790d384c6cb4484",
		"std_msgs/UInt32":              "304a39449588c7f8ce2df6e8001c5fce",
		"std_msgs/UInt32MultiArray":    "4d6a180abc9be191b96a7eda6c8a233d",
		"std_msgs/UInt64":              "1b2a79973e8bf53d7b53acb71299cb57",
		"std_msgs/UInt64MultiArray":    "6088f127afb1d6c72927aa1247e945af",
		"std_msgs/UInt8":               "7c8164229e7d2c17eb95e9231617fdee",
	}

	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")
	ctx, e := NewMsgContext(strings.Split(rosPkgPath, ":"))
	if e != nil {
		t.Errorf("Failed to create MsgContext.")
	} else {
		for fullname, md5 := range std_msgs {
			_, shortName, _ := packageResourceName(fullname)

			t.Run(shortName, func(t *testing.T) {
				var spec *MsgSpec
				spec, e := ctx.LoadMsg(fullname)
				if e != nil {
					t.Errorf("Failed to parse: %v", e)
				} else {
					assertEqual(t, spec.MD5Sum, md5)
				}
			})
		}

	}
}

func TestMD5_sensor_msgs(t *testing.T) {
	var sensor_msgs = map[string]string{
		"sensor_msgs/BatteryState":       "476f837fa6771f6e16e3bf4ef96f8770",
		"sensor_msgs/CameraInfo":         "c9a58c1b0b154e0e6da7578cb991d214",
		"sensor_msgs/ChannelFloat32":     "3d40139cdd33dfedcb71ffeeeb42ae7f",
		"sensor_msgs/CompressedImage":    "8f7a12909da2c9d3332d540a0977563f",
		"sensor_msgs/FluidPressure":      "804dc5cea1c5306d6a2eb80b9833befe",
		"sensor_msgs/Illuminance":        "8cf5febb0952fca9d650c3d11a81a188",
		"sensor_msgs/Image":              "060021388200f6f0f447d0fcd9c64743",
		"sensor_msgs/Imu":                "6a62c6daae103f4ff57a132d6f95cec2",
		"sensor_msgs/JointState":         "3066dcd76a6cfaef579bd0f34173e9fd",
		"sensor_msgs/Joy":                "5a9ea5f83505693b71e785041e67a8bb",
		"sensor_msgs/JoyFeedback":        "f4dcd73460360d98f36e55ee7f2e46f1",
		"sensor_msgs/JoyFeedbackArray":   "cde5730a895b1fc4dee6f91b754b213d",
		"sensor_msgs/LaserEcho":          "8bc5ae449b200fba4d552b4225586696",
		"sensor_msgs/LaserScan":          "90c7ef2dc6895d81024acba2ac42f369",
		"sensor_msgs/MagneticField":      "2f3b0b43eed0c9501de0fa3ff89a45aa",
		"sensor_msgs/MultiDOFJointState": "690f272f0640d2631c305eeb8301e59d",
		"sensor_msgs/MultiEchoLaserScan": "6fefb0c6da89d7c8abe4b339f5c2f8fb",
		"sensor_msgs/NavSatFix":          "2d3a8cd499b9b4a0249fb98fd05cfa48",
		"sensor_msgs/NavSatStatus":       "331cdbddfa4bc96ffc3b9ad98900a54c",
		"sensor_msgs/PointCloud":         "d8e9c3f5afbdd8a130fd1d2763945fca",
		"sensor_msgs/PointCloud2":        "1158d486dd51d683ce2f1be655c3c181",
		"sensor_msgs/PointField":         "268eacb2962780ceac86cbd17e328150",
		"sensor_msgs/Range":              "c005c34273dc426c67a020a87bc24148",
		"sensor_msgs/RegionOfInterest":   "bdb633039d588fcccb441a4d43ccfe09",
		"sensor_msgs/RelativeHumidity":   "8730015b05955b7e992ce29a2678d90f",
		"sensor_msgs/Temperature":        "ff71b307acdbe7c871a5a6d7ed359100",
		"sensor_msgs/TimeReference":      "fded64a0265108ba86c3d38fb11c0c16",
	}

	rosPkgPath := os.Getenv("ROS_PACKAGE_PATH")
	ctx, e := NewMsgContext(strings.Split(rosPkgPath, ":"))
	if e != nil {
		t.Errorf("Failed to create MsgContext.")
	}

	for fullname, md5 := range sensor_msgs {
		_, shortName, _ := packageResourceName(fullname)

		t.Run(shortName, func(t *testing.T) {
			var spec *MsgSpec
			spec, e := ctx.LoadMsg(fullname)
			if e != nil {
				t.Errorf("Failed to parse: %v", e)
			} else {
				assertEqual(t, spec.MD5Sum, md5)
			}
		})
	}
}
