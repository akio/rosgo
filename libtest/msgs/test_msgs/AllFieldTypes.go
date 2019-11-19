// Package test_msgs is automatically generated from the message definition "test_msgs/AllFieldTypes.msg"
package test_msgs

import (
	"bytes"
	"encoding/binary"
	"github.com/edwinhayes/rosgo/ros"
	"github.com/edwinhayes/rosgo/libtest/msgs/std_msgs"
)

const (
	FOO  byte   = 1
	BAR  byte   = 2
	HOGE string = "hoge"
)

type _MsgAllFieldTypes struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgAllFieldTypes) Text() string {
	return t.text
}

func (t *_MsgAllFieldTypes) Name() string {
	return t.name
}

func (t *_MsgAllFieldTypes) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgAllFieldTypes) NewMessage() ros.Message {
	m := new(AllFieldTypes)
	m.H = std_msgs.Header{}
	m.B = 0
	m.I8 = 0
	m.I16 = 0
	m.I32 = 0
	m.I64 = 0
	m.U8 = 0
	m.U16 = 0
	m.U32 = 0
	m.U64 = 0
	m.F32 = 0.0
	m.F64 = 0.0
	m.T = ros.Time{}
	m.D = ros.Duration{}
	m.S = ""
	m.C = std_msgs.ColorRGBA{}
	m.DynAry = []uint32{}
	for i := 0; i < 2; i++ {
		m.FixAry[i] = 0
	}
	return m
}

var (
	MsgAllFieldTypes = &_MsgAllFieldTypes{
		`byte FOO=1
byte BAR=2
string HOGE=hoge

Header h
byte b
int8 i8
int16 i16
int32 i32
int64 i64
uint8 u8
uint16 u16
uint32 u32
uint64 u64
float32 f32
float64 f64
time t
duration d
string s
std_msgs/ColorRGBA c
uint32[] dyn_ary
uint32[2] fix_ary
#std_msgs/ColorRGBA[] msg_ary
`,
		"test_message/AllFieldTypes",
		"5406fac98ad8897d5c798fda29d3f362",
	}
)

type AllFieldTypes struct {
	H      std_msgs.Header    `rosmsg:"h:Header"`
	B      uint8              `rosmsg:"b:byte"`
	I8     int8               `rosmsg:"i8:int8"`
	I16    int16              `rosmsg:"i16:int16"`
	I32    int32              `rosmsg:"i32:int32"`
	I64    int64              `rosmsg:"i64:int64"`
	U8     uint8              `rosmsg:"u8:uint8"`
	U16    uint16             `rosmsg:"u16:uint16"`
	U32    uint32             `rosmsg:"u32:uint32"`
	U64    uint64             `rosmsg:"u64:uint64"`
	F32    float32            `rosmsg:"f32:float32"`
	F64    float64            `rosmsg:"f64:float64"`
	T      ros.Time           `rosmsg:"t:time"`
	D      ros.Duration       `rosmsg:"d:duration"`
	S      string             `rosmsg:"s:string"`
	C      std_msgs.ColorRGBA `rosmsg:"c:ColorRGBA"`
	DynAry []uint32           `rosmsg:"dyn_ary:uint32[]"`
	FixAry [2]uint32          `rosmsg:"fix_ary:uint32[2]"`
}

func (m *AllFieldTypes) Type() ros.MessageType {
	return MsgAllFieldTypes
}

func (m *AllFieldTypes) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	if err = m.H.Serialize(buf); err != nil {
		return err
	}
	binary.Write(buf, binary.LittleEndian, m.B)
	binary.Write(buf, binary.LittleEndian, m.I8)
	binary.Write(buf, binary.LittleEndian, m.I16)
	binary.Write(buf, binary.LittleEndian, m.I32)
	binary.Write(buf, binary.LittleEndian, m.I64)
	binary.Write(buf, binary.LittleEndian, m.U8)
	binary.Write(buf, binary.LittleEndian, m.U16)
	binary.Write(buf, binary.LittleEndian, m.U32)
	binary.Write(buf, binary.LittleEndian, m.U64)
	binary.Write(buf, binary.LittleEndian, m.F32)
	binary.Write(buf, binary.LittleEndian, m.F64)
	binary.Write(buf, binary.LittleEndian, m.T.Sec)
	binary.Write(buf, binary.LittleEndian, m.T.NSec)
	binary.Write(buf, binary.LittleEndian, m.D.Sec)
	binary.Write(buf, binary.LittleEndian, m.D.NSec)
	binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.S))))
	buf.Write([]byte(m.S))
	if err = m.C.Serialize(buf); err != nil {
		return err
	}
	binary.Write(buf, binary.LittleEndian, uint32(len(m.DynAry)))
	for _, e := range m.DynAry {
		binary.Write(buf, binary.LittleEndian, e)
	}
	for _, e := range m.FixAry {
		binary.Write(buf, binary.LittleEndian, e)
	}
	return err
}

func (m *AllFieldTypes) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = m.H.Deserialize(buf); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.B); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.I8); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.I16); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.I32); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.I64); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.U8); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.U16); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.U32); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.U64); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.F32); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.F64); err != nil {
		return err
	}
	{
		if err = binary.Read(buf, binary.LittleEndian, &m.T.Sec); err != nil {
			return err
		}

		if err = binary.Read(buf, binary.LittleEndian, &m.T.NSec); err != nil {
			return err
		}
	}
	{
		if err = binary.Read(buf, binary.LittleEndian, &m.D.Sec); err != nil {
			return err
		}

		if err = binary.Read(buf, binary.LittleEndian, &m.D.NSec); err != nil {
			return err
		}
	}
	{
		var size uint32
		if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
			return err
		}
		data := make([]byte, int(size))
		if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
			return err
		}
		m.S = string(data)
	}
	if err = m.C.Deserialize(buf); err != nil {
		return err
	}
	{
		var size uint32
		if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
			return err
		}
		m.DynAry = make([]uint32, int(size))
		for i := 0; i < int(size); i++ {
			if err = binary.Read(buf, binary.LittleEndian, &m.DynAry[i]); err != nil {
				return err
			}
		}
	}
	{
		for i := 0; i < 2; i++ {
			if err = binary.Read(buf, binary.LittleEndian, &m.FixAry[i]); err != nil {
				return err
			}
		}
	}
	return err
}
