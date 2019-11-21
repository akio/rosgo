// Automatically generated from the message definition "rospy_tutorials/AddTwoInts.srv"
package rospy_tutorials

import (
	"github.com/edwinhayes/rosgo/ros"
)

// Service type metadata
type _SrvAddTwoInts struct {
	name    string
	md5sum  string
	text    string
	reqType ros.MessageType
	resType ros.MessageType
}

func (t *_SrvAddTwoInts) Name() string                  { return t.name }
func (t *_SrvAddTwoInts) MD5Sum() string                { return t.md5sum }
func (t *_SrvAddTwoInts) Text() string                  { return t.text }
func (t *_SrvAddTwoInts) RequestType() ros.MessageType  { return t.reqType }
func (t *_SrvAddTwoInts) ResponseType() ros.MessageType { return t.resType }
func (t *_SrvAddTwoInts) NewService() ros.Service {
	return new(AddTwoInts)
}

var (
	SrvAddTwoInts = &_SrvAddTwoInts{
		"rospy_tutorials/AddTwoInts",
		"6a2e34150c00229791cc89ff309fff21",
		`int64 a
int64 b
---
int64 sum
`,
		MsgAddTwoIntsRequest,
		MsgAddTwoIntsResponse,
	}
)

type AddTwoInts struct {
	Request  AddTwoIntsRequest
	Response AddTwoIntsResponse
}

func (s *AddTwoInts) ReqMessage() ros.Message { return &s.Request }
func (s *AddTwoInts) ResMessage() ros.Message { return &s.Response }
