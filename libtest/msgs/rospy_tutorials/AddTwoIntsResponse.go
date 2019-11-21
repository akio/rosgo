// Automatically generated from the message definition "rospy_tutorials/AddTwoIntsResponse.msg"
package rospy_tutorials

import (
	"bytes"
	"encoding/binary"

	"github.com/edwinhayes/rosgo/ros"
)

type _MsgAddTwoIntsResponse struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgAddTwoIntsResponse) Text() string {
	return t.text
}

func (t *_MsgAddTwoIntsResponse) Name() string {
	return t.name
}

func (t *_MsgAddTwoIntsResponse) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgAddTwoIntsResponse) NewMessage() ros.Message {
	m := new(AddTwoIntsResponse)
	m.Sum = 0
	return m
}

var (
	MsgAddTwoIntsResponse = &_MsgAddTwoIntsResponse{
		`
int64 sum
`,
		"rospy_tutorials/AddTwoIntsResponse",
		"b88405221c77b1878a3cbbfff53428d7",
	}
)

type AddTwoIntsResponse struct {
	Sum int64 `rosmsg:"sum:int64"`
}

func (m *AddTwoIntsResponse) Type() ros.MessageType {
	return MsgAddTwoIntsResponse
}

func (m *AddTwoIntsResponse) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	binary.Write(buf, binary.LittleEndian, m.Sum)
	return err
}

func (m *AddTwoIntsResponse) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = binary.Read(buf, binary.LittleEndian, &m.Sum); err != nil {
		return err
	}
	return err
}
