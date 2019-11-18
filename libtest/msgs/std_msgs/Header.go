// Package std_msgs is automatically generated from the message definition "std_msgs/Header.msg"
package std_msgs

import (
	"bytes"
	"encoding/binary"
	"github.com/edwinhayes/rosgo/ros"
)

type _MsgHeader struct {
	text   string
	name   string
	md5sum string
}

func (t *_MsgHeader) Text() string {
	return t.text
}

func (t *_MsgHeader) Name() string {
	return t.name
}

func (t *_MsgHeader) MD5Sum() string {
	return t.md5sum
}

func (t *_MsgHeader) NewMessage() ros.Message {
	m := new(Header)
	m.Seq = 0
	m.Stamp = ros.Time{}
	m.FrameId = ""
	return m
}

var (
	MsgHeader = &_MsgHeader{
		`# Standard metadata for higher-level stamped data types.
# This is generally used to communicate timestamped data 
# in a particular coordinate frame.
# 
# sequence ID: consecutively increasing ID 
uint32 seq
#Two-integer timestamp that is expressed as:
# * stamp.sec: seconds (stamp_secs) since epoch (in Python the variable is called 'secs')
# * stamp.nsec: nanoseconds since stamp_secs (in Python the variable is called 'nsecs')
# time-handling sugar is provided by the client library
time stamp
#Frame this data is associated with
string frame_id
`,
		"std_msgs/Header",
		"2176decaecbce78abc3b96ef049fabed",
	}
)

type Header struct {
	Seq     uint32   `rosmsg:"seq:uint32"`
	Stamp   ros.Time `rosmsg:"stamp:time"`
	FrameId string   `rosmsg:"frame_id:string"`
}

func (m *Header) Type() ros.MessageType {
	return MsgHeader
}

func (m *Header) Serialize(buf *bytes.Buffer) error {
	var err error = nil
	binary.Write(buf, binary.LittleEndian, m.Seq)
	binary.Write(buf, binary.LittleEndian, m.Stamp.Sec)
	binary.Write(buf, binary.LittleEndian, m.Stamp.NSec)
	binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.FrameId))))
	buf.Write([]byte(m.FrameId))
	return err
}

func (m *Header) Deserialize(buf *bytes.Reader) error {
	var err error = nil
	if err = binary.Read(buf, binary.LittleEndian, &m.Seq); err != nil {
		return err
	}
	{
		if err = binary.Read(buf, binary.LittleEndian, &m.Stamp.Sec); err != nil {
			return err
		}

		if err = binary.Read(buf, binary.LittleEndian, &m.Stamp.NSec); err != nil {
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
		m.FrameId = string(data)
	}
	return err
}
