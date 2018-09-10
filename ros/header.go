// Connection header
package ros

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type header struct {
	key   string
	value string
}

type connHeader struct {
	headers map[string]string
}

const BufferSize = 1024

func readConnectionHeader(r io.Reader) ([]header, error) {
	buf := make([]byte, 4)
	_, err := io.ReadAtLeast(r, buf, 4)
	if err != nil {
		return nil, err
	}
	var headerSize uint32
	bufReader := bytes.NewBuffer(buf)
	err = binary.Read(bufReader, binary.LittleEndian, &headerSize)
	if err != nil {
		return nil, err
	}
	buf = make([]byte, int(headerSize))
	_, err = io.ReadAtLeast(r, buf, int(headerSize))
	if err != nil {
		return nil, err
	}

	var done uint32 = 0
	var headers []header
	bufReader = bytes.NewBuffer(buf)
	for {
		if done == headerSize {
			break
		} else if done > headerSize {
			return nil, fmt.Errorf("Header length overrrun")
		}
		var size uint32
		err := binary.Read(bufReader, binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}
		line := bufReader.Next(int(size))
		sep := bytes.IndexByte(line, '=')
		key := string(line[0:sep])
		value := string(line[sep+1:])
		headers = append(headers, header{key, value})
		done += 4 + size
	}
	return headers, nil
}

func writeConnectionHeader(headers []header, w io.Writer) error {
	//var buf bytes.Buffer
	var headerSize int
	var sizeList []int
	for _, h := range headers {
		size := len(h.key) + len(h.value) + 1
		sizeList = append(sizeList, size)
		headerSize += size + 4
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(headerSize)); err != nil {
		return err
	}
	for i, h := range headers {
		err := binary.Write(w, binary.LittleEndian, uint32(sizeList[i]))
		if err != nil {
			return err
		}
		if _, err = w.Write([]byte(h.key)); err != nil {
			return err
		}
		if _, err = w.Write([]byte("=")); err != nil {
			return err
		}
		if _, err = w.Write([]byte(h.value)); err != nil {
			return err
		}
	}
	return nil
}
