// Simple XMLRPC client/server for go
package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"

	//	"io"
	"net/http"
	//	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func xmlEscape(s string) string {
	var buffer bytes.Buffer
	xml.Escape(&buffer, []byte(s))
	return buffer.String()
}

func emitValue(buf *bytes.Buffer, value interface{}) error {
	if bs, ok := value.([]byte); ok {
		buf.WriteString("<base64>")
		buf.WriteString(base64.StdEncoding.EncodeToString(bs))
		buf.WriteString("</base64>")
	} else {
		val := reflect.ValueOf(value)
		if !val.IsValid() {
			return nil
		}

		t := val.Type()
		k := val.Kind()
		switch k {
		case reflect.Bool:
			b := val.Bool()
			var i int
			if b {
				i = 1
			} else {
				i = 0
			}
			buf.WriteString("<boolean>")
			buf.WriteString(fmt.Sprint(i))
			buf.WriteString("</boolean>")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := val.Int()
			buf.WriteString("<int>")
			buf.WriteString(strconv.FormatInt(i, 10))
			buf.WriteString("</int>")
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u := val.Uint()
			buf.WriteString("<int>")
			buf.WriteString(strconv.FormatInt(int64(u), 10))
			buf.WriteString("</int>")
		case reflect.Float32, reflect.Float64:
			f := val.Float()
			buf.WriteString("<double>")
			buf.WriteString(strconv.FormatFloat(f, 'g', -1, 64))
			buf.WriteString("</double>")
		case reflect.Array, reflect.Slice:
			buf.WriteString("<array><data>")
			for i := 0; i < val.Len(); i++ {
				buf.WriteString("<value>")
				v := val.Index(i)
				e := emitValue(buf, v.Interface())
				if e != nil {
					return e
				}
				buf.WriteString("</value>")
			}
			buf.WriteString("</data></array>")
		case reflect.Map:
			keyKind := t.Key().Kind()
			if keyKind != reflect.String {
				return errors.New("Map key must be string")
			}
			keys := val.MapKeys()
			buf.WriteString("<struct>")
			for _, key := range keys {
				buf.WriteString("<member><name>")
				buf.WriteString(xmlEscape(key.String()))
				buf.WriteString("</name><value>")
				v := val.MapIndex(key)
				e := emitValue(buf, v.Interface())
				if e != nil {
					return e
				}
				buf.WriteString("</value></member>")
			}
			buf.WriteString("</struct>")
		case reflect.String:
			s := val.String()
			buf.WriteString("<string>")
			buf.WriteString(xmlEscape(s))
			buf.WriteString("</string>")
		default:
			return fmt.Errorf("Invalid kind! %v %v", k.String(), val.Type().Name())
		}
	}
	return nil
}

func emitRequest(buf *bytes.Buffer, method string, args ...interface{}) error {
	buf.WriteString(xml.Header)
	buf.WriteString("<methodCall><methodName>")
	buf.WriteString(xmlEscape(method))
	buf.WriteString("</methodName><params>")
	for _, arg := range args {
		buf.WriteString("<param><value>")
		e := emitValue(buf, arg)
		if e != nil {
			return e
		}
		buf.WriteString("</value></param>")
	}
	buf.WriteString("</params></methodCall>")
	return nil
}

func emitResponse(buf *bytes.Buffer, value interface{}) error {
	buf.WriteString(xml.Header)
	buf.WriteString("<methodResponse><params><param><value>")
	e := emitValue(buf, value)
	if e != nil {
		return e
	}
	buf.WriteString("</value></param></params></methodResponse>")
	return nil
}

func emitFault(buf *bytes.Buffer, code int, message string) error {
	buf.WriteString(xml.Header)
	buf.WriteString("<methodResponse><fault><value>")
	fault := make(map[string]interface{})
	fault["faultCode"] = code
	fault["faultString"] = message
	e := emitValue(buf, fault)
	if e != nil {
		return e
	}
	buf.WriteString("</value></fault></methodResponse>")
	return nil
}

func nextTag(d *xml.Decoder) (xml.StartElement, error) {
	for {
		token, e := d.Token()
		if e != nil {
			return xml.StartElement{}, e
		}
		elem, ok := token.(xml.StartElement)
		if ok {
			return elem, nil
		}
	}
	panic("not reached")
}

func expectNextTag(d *xml.Decoder, name string) (xml.StartElement, error) {
	tag, e := nextTag(d)
	if e != nil {
		return xml.StartElement{}, e
	}
	if tag.Name.Local == name {
		return tag, nil
	}
	return xml.StartElement{}, errors.New("Element name mismatch")
}

// Parse a value after the <value> tag has been read.  On (non-error)
// return, the </value> closing tag will have been read.
func parseValue(d *xml.Decoder) (interface{}, error) {
	token, e := d.Token()
	//	t, e := nextTag(d)
	if e != nil {
		return nil, e
	}

	switch t := token.(type) {
	case xml.StartElement:
		switch t.Name.Local {
		case "boolean":
			token, e := d.Token()
			if e != nil {
				return nil, e
			}
			data, ok := token.(xml.CharData)
			if !ok {
				return nil, errors.New("boolean: Not a CharData")
			}
			var i int64
			i, e = strconv.ParseInt(string(data), 10, 4)
			if e != nil {
				return nil, e
			}
			switch i {
			case 0:
				d.Skip() // </bool>
				d.Skip() // </value>
				return false, nil
			case 1:
				d.Skip() // </bool>
				d.Skip() // </value>
				return true, nil
			default:
				return nil, errors.New("Parse error")
			}
		case "i4", "int":
			token, e := d.Token()
			if e != nil {
				return nil, e
			}
			data, ok := token.(xml.CharData)
			if !ok {
				return nil, errors.New("int: Not a CharData")
			}
			var i int64
			i, e = strconv.ParseInt(string(data), 0, 32)
			if e != nil {
				return nil, e
			}
			d.Skip() // </i4> or </int>
			d.Skip() // </value>
			return int32(i), nil
		case "double":
			token, e := d.Token()
			if e != nil {
				return nil, e
			}
			data, ok := token.(xml.CharData)
			if !ok {
				return nil, errors.New("double: Not a CharData")
			}
			var f float64
			f, e = strconv.ParseFloat(string(data), 64)
			if e != nil {
				return nil, e
			}
			d.Skip() // </double>
			d.Skip() // </value>
			return f, nil
		case "string":
			token, e := d.Token()
			if e != nil {
				return nil, e
			}
			data, ok := token.(xml.CharData)
			if ok {
				s := string(data.Copy())
				d.Skip() // </string>
				d.Skip() // </value>
				return s, nil
			} else {
				var end xml.EndElement
				end, ok = token.(xml.EndElement)
				if ok && end.Name.Local == "string" {
					d.Skip() // </value>
					return "", nil
				} else {
					return nil, errors.New("string: parse error")
				}
			}
		case "dateTime.iso8601":
			return nil, errors.New("Not supported1")
		case "base64":
			token, e := d.Token()
			if e != nil {
				return nil, e
			}
			data, ok := token.(xml.CharData)
			if !ok {
				return nil, errors.New("base64: Not a CharData")
			}
			var bs []byte
			bs, e = base64.StdEncoding.DecodeString(string(data))
			if e != nil {
				return nil, e
			}
			d.Skip() // </base64>
			d.Skip() // </value>
			return bs, nil
		case "array":
			_, e := expectNextTag(d, "data")
			if e != nil {
				return nil, e
			}
			var a []interface{}
			for {
				t, e := d.Token()
				if e != nil {
					return nil, e
				}
				switch t.(type) {
				case xml.StartElement:
					elem, _ := t.(xml.StartElement)
					if elem.Name.Local == "value" {
						var val interface{}
						val, e = parseValue(d)
						if e != nil {
							return nil, e
						}
						a = append(a, val)
					}
				case xml.EndElement:
					elem, _ := t.(xml.EndElement)
					if elem.Name.Local == "array" {
						d.Skip() // </value>
						return a, nil
					}
				}
			}
			return nil, errors.New("Not reached")
		case "struct":
			m := make(map[string]interface{})
			var name string
			var value interface{}
			for {
				t, e := d.Token()
				if e != nil {
					return nil, e
				}
				switch t.(type) {
				case xml.StartElement:
					elem, _ := t.(xml.StartElement)
					switch elem.Name.Local {
					case "member":
					case "name":
						t, e = d.Token()
						if e != nil {
							return nil, e
						}
						data, ok := t.(xml.CharData)
						if ok {
							name = string(data)
						} else {
							return nil, errors.New("")
						}
					case "value":
						value, e = parseValue(d)
						if e != nil {
							return nil, e
						}
					}
				case xml.EndElement:
					elem, _ := t.(xml.EndElement)
					switch elem.Name.Local {
					case "member":
						m[name] = value
					case "struct":
						d.Skip() // </value>
						return m, nil
					}
				}
			}
			return nil, errors.New("Not reached")
		default:
			return nil, errors.New("Not supported: t.Name.Local = " + t.Name.Local)
		}
	case xml.CharData:
		copy := t.Copy()
		// spaces and newlines for pretty formatting of xml
		// show up as chardata, so here we ignore them.
		stripped := strings.TrimSpace(string(copy))
		if stripped != "" {
			d.Skip() // </value>
			return string(copy), nil
		} else {
			return parseValue(d)
		}
	case xml.EndElement:
		return "", nil
	}

	return nil, errors.New("Invalid data type")
}

func parseRequest(d *xml.Decoder) (name string, args []interface{}, e error) {
	_, e = expectNextTag(d, "methodCall")
	if e != nil {
		return
	}
	_, e = expectNextTag(d, "methodName")
	if e != nil {
		return
	}
	var t xml.Token
	t, e = d.Token()
	if e != nil {
		return
	}
	data, ok := t.(xml.CharData)
	if !ok {
		e = errors.New("Invalid methodName")
	}
	name = string(data)
	_, e = expectNextTag(d, "params")
	if e != nil {
		return
	}
	for {
		t, e = d.Token()
		switch t.(type) {
		case xml.StartElement:
			elem, _ := t.(xml.StartElement)
			if elem.Name.Local == "value" {
				var x interface{}
				x, e = parseValue(d)
				if e != nil {
					return
				}
				args = append(args, x)
			}
		case xml.EndElement:
			elem, _ := t.(xml.EndElement)
			if elem.Name.Local == "params" {
				d.Skip()
				return
			}
		}
	}
	e = errors.New("Missing end element.")
	return
}

func parseResponse(d *xml.Decoder) (ok bool, result interface{}, e error) {
	_, e = expectNextTag(d, "methodResponse")
	if e != nil {
		return
	}
	var se xml.StartElement
	se, e = nextTag(d)
	if e != nil {
		return
	}
	switch se.Name.Local {
	case "params":
		_, e = expectNextTag(d, "param")
		if e != nil {
			return
		}
		_, e = expectNextTag(d, "value")
		if e != nil {
			return
		}
		result, e = parseValue(d)
		if e != nil {
			return
		}
		ok = true
		d.Skip()
		d.Skip()
		d.Skip()
		return
	case "fault":
		_, e = expectNextTag(d, "value")
		if e != nil {
			return
		}
		result, e = parseValue(d)
		if e != nil {
			return
		}
		ok = false
		d.Skip()
		d.Skip()
		return
	}
	e = errors.New("Missing end element.")
	return
}

// Call a XMLRPC API in a remote host.
// Args:
//   url string: URL of the remote host
func Call(url string, method string, args ...interface{}) (res interface{}, e error) {
	var buffer bytes.Buffer
	e = emitRequest(&buffer, method, args...)
	if e != nil {
		e = fmt.Errorf("Building request failed for %v", e)
		return
	}
	var r *http.Response
	r, e = http.Post(url, "text/xml", &buffer)
	if e != nil {
		e = fmt.Errorf("Sending request failed for %v", e)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("HTTP failed with code %v", r.Status)
		return
	}

	// bodyReader := io.TeeReader(r.Body, os.Stdout)
	// decoder := xml.NewDecoder(bodyReader)
	decoder := xml.NewDecoder(r.Body)
	ok, result, e := parseResponse(decoder)
	if e != nil {
		e = fmt.Errorf("Parsing response failed for %v", e)
		return
	}
	if ok {
		res = result
		return
	} else {
		var m map[string]interface{}
		m, ok = result.(map[string]interface{})
		if ok {
			var c int32
			c, ok = m["faultCode"].(int32)
			if ok {
				var s string
				s, ok = m["faultString"].(string)
				if ok {
					e = fmt.Errorf("XMLRPC Fault: code=%v string=%v", c, s)
					return
				}
			}
		}
		e = errors.New("Malformed XMLRPC Fault Response")
		return
	}
	panic("Not reached")
}

//type Method func (args ...interface{}) (interface{}, error)
type Method interface{}

//
type Handler struct {
	mapping map[string]Method
	wait    sync.WaitGroup
}

//
func NewHandler(mapping map[string]Method) *Handler {
	handler := new(Handler)
	handler.mapping = mapping
	return handler
}

//
func (self *Handler) WaitForShutdown() {
	self.wait.Wait()
}

//
func (self *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.wait.Add(1)
	defer self.wait.Done()

	decoder := xml.NewDecoder(req.Body)
	var buffer bytes.Buffer

	name, args, err := parseRequest(decoder)
	if err != nil {
		err = emitFault(&buffer, 1, "Invalid request.")
		_, err = buffer.WriteTo(w)
		return
	}

	method, ok := self.mapping[name]
	if !ok {
		err = emitFault(&buffer, 1, fmt.Sprintf("No method named '%v'.", name))
		_, err = buffer.WriteTo(w)
		return
	}

	argValues := []reflect.Value{}
	for _, v := range args {
		argValues = append(argValues, reflect.ValueOf(v))
	}
	resultValues := reflect.ValueOf(method).Call(argValues)
	if len(resultValues) != 2 {
		err = emitFault(&buffer, 1, fmt.Sprintf("Method '%v' return invalid results.", name))
		return
	}
	errValue := resultValues[1]
	if !errValue.IsNil() {
		err, ok = errValue.Interface().(error)
		if !ok {
			err = emitFault(&buffer, 1, fmt.Sprintf("Method '%v' return an invalid error.", name))
			_, err = buffer.WriteTo(w)
		} else {
			err = emitFault(&buffer, 1, fmt.Sprintf("Method '%v' call failed.", name))
			_, err = buffer.WriteTo(w)
		}
		return
	}

	err = emitResponse(&buffer, resultValues[0].Interface())
	if err != nil {
		err = emitFault(&buffer, 1, fmt.Sprintf("Method '%v' return an invalid result type.", name))
		_, err = buffer.WriteTo(w)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
	_, err = buffer.WriteTo(w)
	w.(http.Flusher).Flush()
}
