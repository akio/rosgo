// Automatically generated from the message definition "std_msgs/String.msg"
package std_msgs

import (
	"bytes"
	"encoding/binary"

	"github.com/edwinhayes/rosgo/ros"
)

type _MsgString struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgString) Text() string {
	return t.text
}

func (t *_MsgString) Name() string {
	return t.name
}

func (t *_MsgString) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgString) NewMessage() ros.Message {
	m := new(String)
	m.Data = ""
	return m
}

var (
	MsgString = &_MsgString{
		`string data
`,
		"std_msgs/String",
		"992ce8a1687cec8c8bd883ec73ca41d1",
	}
)

type String struct {
	Data string `rosmsg:"data:string"`
}

func (m *String) Type() ros.MessageType {
	return MsgString
}

func (m *String) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.Data))))
	buf.Write([]byte(m.Data))
	return err
}

func (m *String) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	{
		var size uint32
		if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
			return err
		}
		data := make([]byte, int(size))
		if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
			return err
		}
		m.Data = string(data)
	}
	return err
}
