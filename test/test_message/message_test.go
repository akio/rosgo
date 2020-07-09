package test_message

//go:generate gengo msg rosgo_tests/AllFieldTypes AllFieldTypes.msg
//go:generate gengo msg std_msgs/Header
//go:generate gengo msg std_msgs/Int16
//go:generate gengo msg std_msgs/Int32
//go:generate gengo msg std_msgs/ColorRGBA
import (
	"bytes"
	"fmt"
	"rosgo_tests"
	"std_msgs"
	"testing"

	"github.com/fetchrobotics/rosgo/ros"
)

func TestInitialize(t *testing.T) {
	msg := rosgo_tests.AllFieldTypes{}
	fmt.Println(msg)
	fmt.Println(msg.H)

	if msg.B != 0 {
		t.Error(msg.B)
	}

	if msg.I8 != 0 {
		t.Error(msg.I8)
	}

	if msg.I16 != 0 {
		t.Error(msg.I16)
	}

	if msg.I32 != 0 {
		t.Error(msg.I32)
	}

	if msg.I64 != 0 {
		t.Error(msg.I64)
	}

	if msg.U8 != 0 {
		t.Error(msg.U8)
	}

	if msg.U16 != 0 {
		t.Error(msg.U16)
	}

	if msg.U32 != 0 {
		t.Error(msg.U32)
	}

	if msg.U64 != 0 {
		t.Error(msg.U64)
	}

	if msg.F32 != 0.0 {
		t.Error(msg.F32)
	}

	if msg.F64 != 0.0 {
		t.Error(msg.F64)
	}

	if msg.T.Sec != 0 || msg.T.NSec != 0 {
		t.Error(msg.T)
	}

	if msg.D.Sec != 0 || msg.D.NSec != 0 {
		t.Error(msg.D)
	}

	if msg.S != "" {
		t.Error(msg.S)
	}

	if msg.C.R != 0.0 || msg.C.G != 0.0 || msg.C.B != 0.0 || msg.C.A != 0 {
		t.Error(msg.C)
	}

	if len(msg.DynAry) != 0 {
		t.Error(msg.DynAry)
	}

	if len(msg.FixAry) != 2 || msg.FixAry[0] != 0 || msg.FixAry[1] != 0 {
		t.Error(msg.FixAry)
	}

}

func CheckBytes(t *testing.T, a, b []byte) {
	if !bytes.Equal(a, b) {
		if len(a) != len(b) {
			t.Errorf("expected length is %d but %d", len(a), len(b))
		} else {
			for i := 0; i < len(a); i++ {
				if a[i] != b[i] {
					t.Errorf("result[%3d] is expected to be %02X but %02X", i, a[i], b[i])
				} else {
					t.Errorf("%02X", a[i])
				}
			}
		}
	}
}

func TestSerializeHeader(t *testing.T) {
	var msg std_msgs.Header
	msg.Seq = 0x89ABCDEF
	msg.Stamp = ros.NewTime(0x89ABCDEF, 0x01234567)
	msg.FrameId = "frame_id"
	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}
	result := buf.Bytes()
	expected := []byte{
		// Header.Seq
		0xEF, 0xCD, 0xAB, 0x89,
		// Header.Stamp
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// Header.FrameId
		0x08, 0x00, 0x00, 0x00,
		0x66, 0x72, 0x61, 0x6D, 0x65, 0x5F, 0x69, 0x64,
	}
	CheckBytes(t, expected, result)
}

func TestSerializeInt16(t *testing.T) {
	var msg std_msgs.Int16
	msg.Data = 0x0123
	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}
	result := buf.Bytes()
	expected := []byte{
		0x23, 0x01,
	}
	CheckBytes(t, expected, result)
}

func TestSerializeInt32(t *testing.T) {
	var msg std_msgs.Int32
	msg.Data = 0x01234567
	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}
	result := buf.Bytes()
	expected := []byte{
		0x67, 0x45, 0x23, 0x01,
	}
	CheckBytes(t, expected, result)
}

func getTestData() []byte {
	return []byte{
		// Header.Seq
		0xEF, 0xCD, 0xAB, 0x89,
		// Header.Stamp
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// Header.FrameId
		0x08, 0x00, 0x00, 0x00,
		0x66, 0x72, 0x61, 0x6D, 0x65, 0x5F, 0x69, 0x64,
		// B
		0x01,
		// I8
		0x01,
		// I16
		0x23, 0x01,
		// I32
		0x67, 0x45, 0x23, 0x01,
		// I64
		0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01,
		// U8
		0x01,
		// U16
		0x23, 0x01,
		// U32
		0x67, 0x45, 0x23, 0x01,
		// U64
		0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01,
		// F32
		0xDB, 0x0F, 0x49, 0x40,
		// F64
		0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40,
		// T
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// D
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// S
		0x0D, 0x00, 0x00, 0x00,
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x21,
		// C
		0x00, 0x00, 0x80, 0x3F,
		0x00, 0x00, 0x00, 0x3F,
		0x00, 0x00, 0x80, 0x3E,
		0x00, 0x00, 0x00, 0x3E,
		// DynAry
		0x02, 0x00, 0x00, 0x00,
		0x67, 0x45, 0x23, 0x01,
		0xEF, 0xCD, 0xAB, 0x89,
		// FixAry
		0x67, 0x45, 0x23, 0x01,
		0xEF, 0xCD, 0xAB, 0x89,
	}
}

func TestSerialize(t *testing.T) {
	var msg rosgo_tests.AllFieldTypes

	msg.H.Seq = 0x89ABCDEF
	msg.H.Stamp = ros.NewTime(0x89ABCDEF, 0x01234567)
	msg.H.FrameId = "frame_id"
	msg.B = 0x01
	msg.I8 = 0x01
	msg.I16 = 0x0123
	msg.I32 = 0x01234567
	msg.I64 = 0x0123456789ABCDEF
	msg.U8 = 0x01
	msg.U16 = 0x0123
	msg.U32 = 0x01234567
	msg.U64 = 0x0123456789ABCDEF
	msg.F32 = 3.141592653589793238462643383
	msg.F64 = 3.1415926535897932384626433832795028842
	msg.T = ros.NewTime(0x89ABCDEF, 0x01234567)
	msg.D = ros.NewDuration(0x89ABCDEF, 0x01234567)
	msg.S = "Hello, world!"
	msg.C = std_msgs.ColorRGBA{1.0, 0.5, 0.25, 0.125}

	msg.DynAry = append(msg.DynAry, 0x01234567)
	msg.DynAry = append(msg.DynAry, 0x89ABCDEF)
	msg.FixAry[0] = 0x01234567
	msg.FixAry[1] = 0x89ABCDEF

	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}

	result := buf.Bytes()
	expected := getTestData()
	CheckBytes(t, expected, result)
}

func TestDeserialize(t *testing.T) {
	source := getTestData()
	reader := bytes.NewReader(source)
	var msg rosgo_tests.AllFieldTypes
	err := msg.Deserialize(reader)
	if err != nil {
		t.Error(err)
	}

	if msg.H.Seq != 0x89ABCDEF {
		t.Error(msg.H.Seq)
	}
	if msg.H.Stamp.Sec != 0x89ABCDEF || msg.H.Stamp.NSec != 0x01234567 {
		t.Error(msg.H.Stamp)
	}
	if msg.H.FrameId != "frame_id" {
		t.Error(msg.H.FrameId)
	}
	if msg.B != 0x01 {
		t.Error(msg.B)
	}
	if msg.I8 != 0x01 {
		t.Error(msg.I8)
	}
	if msg.I16 != 0x0123 {
		t.Error(msg.I16)
	}
	if msg.I32 != 0x01234567 {
		t.Error(msg.I32)
	}
	if msg.I64 != 0x0123456789ABCDEF {
		t.Error(msg.I64)
	}
	if msg.U8 != 0x01 {
		t.Error(msg.U8)
	}
	if msg.U16 != 0x0123 {
		t.Error(msg.U16)
	}
	if msg.U32 != 0x01234567 {
		t.Error(msg.U32)
	}
	if msg.U64 != 0x0123456789ABCDEF {
		t.Error(msg.U64)
	}
	if msg.F32 != 3.141592653589793238462643383 {
		t.Error(msg.F32)
	}
	if msg.F64 != 3.1415926535897932384626433832795028842 {
		t.Error(msg.F64)
	}
	if msg.T.Sec != 0x89ABCDEF || msg.T.NSec != 0x01234567 {
		t.Error(msg.T)
	}
	if msg.D.Sec != 0x89ABCDEF || msg.D.NSec != 0x01234567 {
		t.Error(msg.D)
	}
	if msg.S != "Hello, world!" {
		t.Error(msg.S)
	}
	if msg.C.R != 1.0 || msg.C.G != 0.5 || msg.C.B != 0.25 || msg.C.A != 0.125 {
		t.Error(msg.C)
	}
	if msg.DynAry[0] != 0x01234567 || msg.DynAry[1] != 0x89ABCDEF {
		t.Error(msg.DynAry)
	}
	if msg.FixAry[0] != 0x01234567 || msg.FixAry[1] != 0x89ABCDEF {
		t.Error(msg.DynAry)
	}
	if reader.Len() != 0 {
		t.Fail()
	}
}
