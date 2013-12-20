package rosgo_test

import (
    "bytes"
    "encoding/binary"
    "ros"
)

// Message type metadata
type type_AddTwoIntsRequest struct {
    definition string
    name       string
    md5sum     string
}

func (t *type_AddTwoIntsRequest) Definition() string      { return t.definition }
func (t *type_AddTwoIntsRequest) Name() string            { return t.name }
func (t *type_AddTwoIntsRequest) MD5Sum() string          { return t.md5sum }
func (t *type_AddTwoIntsRequest) NewMessage() ros.Message {
    return new(AddTwoIntsRequest)
}

var (
    TypeOfAddTwoIntsRequest = &type_AddTwoIntsRequest{
        "",
        "rosgo_test/AddTwoIntsRequest",
        "ef8322123148e475e3e93a1f609b2f70",
    }
)

type AddTwoIntsRequest struct {
    A int32
    B int32
}

func (m *AddTwoIntsRequest) Serialize(buf *bytes.Buffer) error {
    var buf bytes.Buffer
    binary.Write(&buf, binary.LittleEndian, m.A)
    binary.Write(&buf, binary.LittleEndian, m.B)
    return nil
}

func (m *AddTwoIntsRequest) Deserialize(buf *bytes.Reader) error {
    if err := binary.Read(buf, binary.LittleEndian, &m.A); err != nil {
        return err
    }
    if err := binary.Read(buf, binary.LittleEndian, &m.B); err != nil {
        return err
    }
    return nil
}


// Message type metadata
type type_AddTwoIntsResponse struct {
    definition string
    name       string
    md5sum     string
}

func (t *type_AddTwoIntsResponse) Definition() string      { return t.definition }
func (t *type_AddTwoIntsResponse) Name() string            { return t.name }
func (t *type_AddTwoIntsResponse) MD5Sum() string          { return t.md5sum }
func (t *type_AddTwoIntsResponse) NewMessage() ros.Message {
    return new(AddTwoIntsRequest)
}

var (
    TypeOfAddTwoIntsResponse = &type_AddTwoIntsResponse{
        "",
        "rosgo_test/AddTwoIntsResponse",
        "0ba699c25c9418c0366f3595c0c8e8ec",
    }
)

type AddTwoIntsResponse struct {
    Sum int32
}

func (m *AddTwoIntsResponse) Serialize(buf *bytes.Buffer) error {
    binary.Write(&buf, binary.LittleEndian, m.Sum)
    return nil
}

func (m *AddTwoIntsResponse) Deserialize(buf *bytes.Reader) error {
    if err := binary.Read(buf, binary.LittleEndian, &m.Sum); err != nil {
        return err
    }
    return nil
}

// Service type metadata
type type_AddTwoInts struct {
    name string
    md5sum string
    reqType ros.MessageType
    resType ros.MessageType
}

func (t *type_AddTwoInts) Name() string { return t.name }
func (t *type_AddTwoInts) MD5Sum() string { return t.md5sum }
func (t *type_AddTwoInts) RequestType() ros.MessageType { return t.reqType }
func (t *type_AddTwoInts) ResponseType() ros.MessageType { return t.resType }
func (t *type_AddTwoInts) NewService() ros.Service {
    return new(AddTwoInts)
}

var (
    TypeOfAddTwoInts = &type_AddTwoInts {
        "rosgo_test/AddTwoInts",
        "f0b6d69ea10b0cf210cb349d58d59e8f",
        TypeOfAddTwoIntsRequest,
        TypeOfAddTwoIntsResponse,
    }
)


type AddTwoInts struct {
    Request AddTwoIntsRequest
    Response AddTwoIntsResponse
}

func (s *AddTwoInts) ReqMessage() ros.Message { return &s.Request }
func (s *AddTwoInts) ResMessage() ros.Message { return &s.Response }

