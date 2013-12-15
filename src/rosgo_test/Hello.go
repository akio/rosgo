package rosgo_test

import (
    "bytes"
    "encoding/binary"
    "ros"
)

type _HelloType struct {
    definition string
    name       string
    md5sum     string
}

func (t *_HelloType) Definition() string      { return t.definition }
func (t *_HelloType) Name() string            { return t.name }
func (t *_HelloType) MD5Sum() string          { return t.md5sum }
func (t *_HelloType) NewMessage() ros.Message { return new(Hello) }

func TypeOfHello() ros.MessageType {
    t := _HelloType{
        "",
        "rosgo_test/Hello",
        "992ce8a1687cec8c8bd883ec73ca41d1",
    }
    return &t
}

type Hello struct {
    Data string
}

func (s *Hello) Serialize() []byte {
    var buf bytes.Buffer
    data := []byte(s.Data)
    size := uint32(len(data))
    binary.Write(&buf, binary.LittleEndian, size)
    buf.Write(data)
    return buf.Bytes()
}

func (s *Hello) Deserialize(buffer []byte) error {
    buf := bytes.NewBuffer(buffer)
    var size uint32
    if err := binary.Read(buf, binary.LittleEndian, &size); err != nil {
        return err
    }
    data := make([]byte, int(size))
    if _, err := buf.Read(data); err != nil {
        return err
    }
    s.Data = string(data)
    return nil
}
