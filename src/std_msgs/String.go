package std_msgs

import (
    "bytes"
    "encoding/binary"
    "ros"
)

type _StringType struct {
    definition string
    name       string
    md5sum     string
}

func (t *_StringType) Definition() string      { return t.definition }
func (t *_StringType) Name() string            { return t.name }
func (t *_StringType) MD5Sum() string          { return t.md5sum }
func (t *_StringType) NewMessage() ros.Message { return new(String) }

func TypeOfString() ros.MessageType {
    t := _StringType{
        "",
        "std_msgs/String",
        "992ce8a1687cec8c8bd883ec73ca41d1",
    }
    return &t
}

type String struct {
    Data string
}

func (s *String) Serialize() []byte {
    var buf bytes.Buffer
    data := []byte(s.Data)
    binary.Write(&buf, binary.LittleEndian, len(data))
    buf.Write(data)
    return buf.Bytes()
}

func (s *String) Deserialize(buffer []byte) error {
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
