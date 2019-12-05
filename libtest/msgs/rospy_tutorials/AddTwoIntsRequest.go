// Automatically generated from the message definition "rospy_tutorials/AddTwoIntsRequest.msg"
package rospy_tutorials

import (
	"bytes"
	"encoding/binary"

	"github.com/edwinhayes/rosgo/ros"
)

type _MsgAddTwoIntsRequest struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgAddTwoIntsRequest) Text() string {
	return t.text
}

func (t *_MsgAddTwoIntsRequest) Name() string {
	return t.name
}

func (t *_MsgAddTwoIntsRequest) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgAddTwoIntsRequest) NewMessage() ros.Message {
	m := new(AddTwoIntsRequest)
	m.A = 0
	m.B = 0
	return m
}

var (
	MsgAddTwoIntsRequest = &_MsgAddTwoIntsRequest{
		`int64 a
int64 b
`,
		"rospy_tutorials/AddTwoIntsRequest",
		"36d09b846be0b371c5f190354dd3153e",
	}
)

type AddTwoIntsRequest struct {
	A int64 `rosmsg:"a:int64"`
	B int64 `rosmsg:"b:int64"`
}

func (m *AddTwoIntsRequest) Type() ros.MessageType {
	return MsgAddTwoIntsRequest
}

func (m *AddTwoIntsRequest) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	binary.Write(buf, binary.LittleEndian, m.A)
	binary.Write(buf, binary.LittleEndian, m.B)
	return err
}

func (m *AddTwoIntsRequest) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = binary.Read(buf, binary.LittleEndian, &m.A); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &m.B); err != nil {
		return err
	}
	return err
}
