// Package std_msgs is automatically generated from the message definition "std_msgs/ColorRGBA.msg"
package std_msgs

import (
	"bytes"
	"encoding/binary"
	"github.com/edwinhayes/rosgo/ros"
)

type _MsgColorRGBA struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgColorRGBA) Text() string {
	return t.text
}

func (t *_MsgColorRGBA) Name() string {
	return t.name
}

func (t *_MsgColorRGBA) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgColorRGBA) NewMessage() ros.Message {
	m := new(ColorRGBA)
	m.R = 0.0
	m.G = 0.0
	m.B = 0.0
	m.A = 0.0
	return m
}

var (
	MsgColorRGBA = &_MsgColorRGBA{
		`float32 r
float32 g
float32 b
float32 a
`,
		"std_msgs/ColorRGBA",
		"a29a96539573343b1310c73607334b00",
	}
)

type ColorRGBA struct {
	R float32 `rosmsg:"r:float32"`
	G float32 `rosmsg:"g:float32"`
	B float32 `rosmsg:"b:float32"`
	A float32 `rosmsg:"a:float32"`
}

func (m *ColorRGBA) Type() ros.MessageType {
	return MsgColorRGBA
}

func (m *ColorRGBA) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	binary.Write(buf, binary.LittleEndian, m.R)
	binary.Write(buf, binary.LittleEndian, m.G)
	binary.Write(buf, binary.LittleEndian, m.B)
	binary.Write(buf, binary.LittleEndian, m.A)
	return err
}

func (m *ColorRGBA) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = binary.Read(buf, binary.LittleEndian, &m.R); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.G); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.B); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.A); err != nil {
		return err
	}
	return err
}
