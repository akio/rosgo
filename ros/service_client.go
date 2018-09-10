package ros

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"
)

type defaultServiceClient struct {
	logger    Logger
	service   string
	srvType   ServiceType
	masterUri string
	nodeId    string
}

func newDefaultServiceClient(logger Logger, nodeId string, masterUri string, service string, srvType ServiceType) *defaultServiceClient {
	client := new(defaultServiceClient)
	client.logger = logger
	client.service = service
	client.srvType = srvType
	client.masterUri = masterUri
	client.nodeId = nodeId
	return client
}

func (c *defaultServiceClient) Call(srv Service) error {
	logger := c.logger

	result, err := callRosApi(c.masterUri, "lookupService", c.nodeId, c.service)
	if err != nil {
		return err
	}

	serviceRawUrl, converted := result.(string)
	if !converted {
		return fmt.Errorf("Result of 'lookupService' is not a string")
	}
	var serviceUrl *url.URL
	serviceUrl, err = url.Parse(serviceRawUrl)
	if err != nil {
		return err
	}

	var conn net.Conn
	conn, err = net.Dial("tcp", serviceUrl.Host)
	if err != nil {
		return err
	}

	// 1. Write connection header
	var headers []header
	md5sum := c.srvType.MD5Sum()
	msgType := c.srvType.Name()
	headers = append(headers, header{"service", c.service})
	headers = append(headers, header{"md5sum", md5sum})
	headers = append(headers, header{"type", msgType})
	headers = append(headers, header{"callerid", c.nodeId})
	logger.Debug("TCPROS Connection Header")
	for _, h := range headers {
		logger.Debugf("  `%s` = `%s`", h.key, h.value)
	}
	conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
	if err := writeConnectionHeader(headers, conn); err != nil {
		return err
	}

	// 2. Read reponse header
	conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
	if resHeaders, err := readConnectionHeader(conn); err != nil {
		return err
	} else {
		logger.Debug("TCPROS Response Header:")
		resHeaderMap := make(map[string]string)
		for _, h := range resHeaders {
			resHeaderMap[h.key] = h.value
			logger.Debugf("  `%s` = `%s`", h.key, h.value)
		}
		if resHeaderMap["type"] != msgType || resHeaderMap["md5sum"] != md5sum {
			logger.Fatalf("Incompatible message type!")
		}
		logger.Debug("Start receiving messages...")
	}

	// 3. Send request
	var buf bytes.Buffer
	_ = srv.ReqMessage().Serialize(&buf)
	reqMsg := buf.Bytes()
	size := uint32(len(reqMsg))
	conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
	if err := binary.Write(conn, binary.LittleEndian, size); err != nil {
		return err
	}
	logger.Debug(len(reqMsg))
	conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
	if _, err := conn.Write(reqMsg); err != nil {
		return err
	}

	// 4. Read OK byte
	var ok byte
	conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
	if err := binary.Read(conn, binary.LittleEndian, &ok); err != nil {
		return err
	} else {
		if ok == 0 {
			var size uint32
			conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
			if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
				return err
			} else {
				errMsg := make([]byte, int(size))
				conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
				if _, err := io.ReadFull(conn, errMsg); err != nil {
					return err
				} else {
					return errors.New(string(errMsg))
				}
			}
		}
	}

	// 5. Receive response
	conn.SetDeadline(time.Now().Add(10 * time.Millisecond))
	//logger.Debug("Reading message size...")
	var msgSize uint32
	if err := binary.Read(conn, binary.LittleEndian, &msgSize); err != nil {
		return err
	}
	logger.Debugf("  %d", msgSize)
	resBuffer := make([]byte, int(msgSize))
	//logger.Debug("Reading message body...")
	if _, err = io.ReadFull(conn, resBuffer); err != nil {
		return err
	}
	resReader := bytes.NewReader(resBuffer)
	if err := srv.ResMessage().Deserialize(resReader); err != nil {
		return err
	}
	return nil
}

func (*defaultServiceClient) Shutdown() {}
