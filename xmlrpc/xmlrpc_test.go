package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestEmitNil(t *testing.T) {
	var buffer bytes.Buffer
	e := emitValue(&buffer, nil)
	if e != nil {
		t.Error(e)
	}
}

func TestEmitBoolean(t *testing.T) {
	var buffer bytes.Buffer
	e := emitValue(&buffer, true)
	if e != nil {
		t.Error(e)
	}
	e = emitValue(&buffer, false)
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	if s != "<boolean>1</boolean><boolean>0</boolean>" {
		t.Error(s)
	}
}

func TestEmitInt(t *testing.T) {
	var buffer bytes.Buffer
	e := emitValue(&buffer, 42)
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	if s != "<int>42</int>" {
		t.Error(s)
	}
}

func TestEmitDouble(t *testing.T) {
	var buffer bytes.Buffer
	e := emitValue(&buffer, 3.14)
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	if s != "<double>3.14</double>" {
		t.Error(s)
	}
}

func TestEmitString(t *testing.T) {
	var buffer bytes.Buffer
	e := emitValue(&buffer, "Hello, world!")
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	if s != "<string>Hello, world!</string>" {
		t.Error(s)
	}
}

func TestEmitBase64(t *testing.T) {
	var buffer bytes.Buffer
	e := emitValue(&buffer, []byte("ABCDEFG"))
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	if s != "<base64>QUJDREVGRw==</base64>" {
		t.Error(s)
	}
}

func TestEmitArray(t *testing.T) {
	var buffer bytes.Buffer
	xs := [...]interface{}{12, "Egypt", false, -31}
	e := emitValue(&buffer, xs)
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	expected := "<array><data>"
	expected += "<value><int>12</int></value>"
	expected += "<value><string>Egypt</string></value>"
	expected += "<value><boolean>0</boolean></value>"
	expected += "<value><int>-31</int></value>"
	expected += "</data></array>"
	if s != expected {
		t.Error(s)
	}
}

func TestEmitArrayFromSlice(t *testing.T) {
	var buffer bytes.Buffer
	xs := []interface{}{12, "Egypt", false, -31}
	e := emitValue(&buffer, xs)
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	expected := "<array><data>"
	expected += "<value><int>12</int></value>"
	expected += "<value><string>Egypt</string></value>"
	expected += "<value><boolean>0</boolean></value>"
	expected += "<value><int>-31</int></value>"
	expected += "</data></array>"
	if s != expected {
		t.Error(s)
	}
}

func TestEmitStruct(t *testing.T) {
	var buffer bytes.Buffer
	xs := make(map[string]interface{})
	xs["lowerBound"] = 18
	xs["upperBound"] = 139
	e := emitValue(&buffer, xs)
	if e != nil {
		t.Error(e)
	}
	s := buffer.String()
	expected1 := "<struct><member>"
	expected1 += "<name>lowerBound</name>"
	expected1 += "<value><int>18</int></value>"
	expected1 += "</member><member>"
	expected1 += "<name>upperBound</name>"
	expected1 += "<value><int>139</int></value>"
	expected1 += "</member></struct>"

	expected2 := "<struct><member>"
	expected2 += "<name>upperBound</name>"
	expected2 += "<value><int>139</int></value>"
	expected2 += "</member><member>"
	expected2 += "<name>lowerBound</name>"
	expected2 += "<value><int>18</int></value>"
	expected2 += "</member></struct>"
	if s != expected1 && s != expected2 {
		t.Error(s)
	}
}

func TestEmitRequest(t *testing.T) {
	var buffer bytes.Buffer
	emitRequest(&buffer, "doSomething", true, 42)
	s := buffer.String()
	expected := xml.Header
	expected += "<methodCall>"
	expected += "<methodName>doSomething</methodName>"
	expected += "<params>"
	expected += "<param><value><boolean>1</boolean></value></param>"
	expected += "<param><value><int>42</int></value></param>"
	expected += "</params>"
	expected += "</methodCall>"
	if s != expected {
		t.Error(s)
	}
}

func TestEmitResponse(t *testing.T) {
	var buffer bytes.Buffer
	emitResponse(&buffer, 42)
	s := buffer.String()
	expected := xml.Header
	expected += "<methodResponse>"
	expected += "<params><param>"
	expected += "<value><int>42</int></value>"
	expected += "</param></params>"
	expected += "</methodResponse>"
	if s != expected {
		t.Error(s)
	}
}

func TestEmitFault(t *testing.T) {
	var buffer bytes.Buffer
	emitFault(&buffer, 42, "failed")
	s := buffer.String()
	expected1 := xml.Header
	expected1 += "<methodResponse><fault><value>"
	expected1 += "<struct><member>"
	expected1 += "<name>faultCode</name>"
	expected1 += "<value><int>42</int></value>"
	expected1 += "</member><member>"
	expected1 += "<name>faultString</name>"
	expected1 += "<value><string>failed</string></value>"
	expected1 += "</member></struct>"
	expected1 += "</value></fault></methodResponse>"

	expected2 := xml.Header
	expected2 += "<methodResponse><fault><value>"
	expected2 += "<struct><member>"
	expected2 += "<name>faultString</name>"
	expected2 += "<value><string>failed</string></value>"
	expected2 += "</member><member>"
	expected2 += "<name>faultCode</name>"
	expected2 += "<value><int>42</int></value>"
	expected2 += "</member></struct>"
	expected2 += "</value></fault></methodResponse>"
	if s != expected1 && s != expected2 {
		t.Error(s)
	}
}

func TestParseBoolean(t *testing.T) {
	buffer := bytes.NewBufferString("<value><boolean>0</boolean></value><value><boolean>1</boolean></value>")
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.(bool)
	if !ok {
		t.Error(ok)
	}
	if x {
		t.Error(x)
	}

	_, _ = decoder.Token() // <value>
	value, e = parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok = value.(bool)
	if !ok {
		t.Error(ok)
	}
	if !x {
		t.Error(x)
	}
}

func TestParseInt(t *testing.T) {
	buffer := bytes.NewBufferString("<value><int>-432</int></value><value><i4>43</i4></value>")
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.(int32)
	if !ok {
		t.Error(ok)
	}
	if x != -432 {
		t.Error(x)
	}

	_, _ = decoder.Token() // <value>
	value, e = parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok = value.(int32)
	if !ok {
		t.Error(ok)
	}
	if x != 43 {
		t.Error(x)
	}
}

func TestParseDouble(t *testing.T) {
	buffer := bytes.NewBufferString("<value><double>-273.5</double></value><value><double>3.14</double></value>")
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.(float64)
	if !ok {
		t.Error(ok)
	}
	if x != -273.5 {
		t.Error(x)
	}

	_, _ = decoder.Token() // <value>
	value, e = parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok = value.(float64)
	if !ok {
		t.Error(ok)
	}
	if x != 3.14 {
		t.Error(x)
	}
}

func TestParseString(t *testing.T) {
	buffer := bytes.NewBufferString("<value><string>Hello, world!</string></value>")
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.(string)
	if !ok {
		t.Error(ok)
	}
	if x != "Hello, world!" {
		t.Error(x)
	}
}

func TestParseBase64(t *testing.T) {
	buffer := bytes.NewBufferString("<value><base64>QUJDREVGRw==</base64></value>")
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.([]byte)
	if !ok {
		t.Error(ok)
	}
	if string(x) != "ABCDEFG" {
		t.Error(x)
	}
}

func TestParseArray(t *testing.T) {
	source := `<value><array>
                   <data>
                       <value><i4>12</i4></value>
                       <value><string>Egypt</string></value>
                       <value><boolean>0</boolean></value>
                       <value><i4>-31</i4></value>
                   </data>
               </array></value>`
	buffer := bytes.NewBufferString(source)
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.([]interface{})
	if !ok {
		t.Error(ok)
	}
	if len(x) != 4 {
		t.Error(x)
	}

	var i int32
	i, ok = x[0].(int32)
	if !ok {
		t.Error(ok)
	}
	if i != 12 {
		t.Error(i)
	}

	var s string
	s, ok = x[1].(string)
	if !ok {
		t.Error(ok)
	}
	if s != "Egypt" {
		t.Error(s)
	}

	var b bool
	b, ok = x[2].(bool)
	if !ok {
		t.Error(ok)
	}
	if b != false {
		t.Error(b)
	}

	i, ok = x[3].(int32)
	if !ok {
		t.Error(ok)
	}
	if i != -31 {
		t.Error(i)
	}
}

func TestParseStruct(t *testing.T) {
	source := `<value><struct>
                   <member>
                       <name>lowerBound</name>
                       <value><i4>18</i4></value>
                   </member>
                   <member>
                       <name>upperBound</name>
                       <value><i4>139</i4></value>
                   </member>
               </struct></value>`
	buffer := bytes.NewBufferString(source)
	decoder := xml.NewDecoder(buffer)
	_, _ = decoder.Token() // <value>
	value, e := parseValue(decoder)
	if e != nil {
		t.Error(e)
	}
	x, ok := value.(map[string]interface{})
	if !ok {
		t.Error(ok)
	}
	if len(x) != 2 {
		t.Error()
	}

	var i int32
	i, ok = x["lowerBound"].(int32)
	if !ok {
		t.Error(ok)
	}
	if i != 18 {
		t.Error(i)
	}

	i, ok = x["upperBound"].(int32)
	if !ok {
		t.Error(ok)
	}
	if i != 139 {
		t.Error(i)
	}
}

func TestParseRequest(t *testing.T) {
	source := xml.Header
	source += `<methodCall>
                   <methodName>doSomething</methodName>
                   <params>
                       <param><value><boolean>1</boolean></value></param>
                       <param><value><int>42</int></value></param>
                   </params>
               </methodCall>`
	buffer := bytes.NewBufferString(source)
	decoder := xml.NewDecoder(buffer)
	name, args, e := parseRequest(decoder)
	if e != nil {
		t.Error(e)
		return
	}
	if name != "doSomething" {
		t.Error(name)
		return
	}
	if len(args) != 2 {
		t.Error(args)
		return
	}

	b, ok := args[0].(bool)
	if !ok {
		t.Error(ok)
		return
	}
	if b != true {
		t.Error(b)
		return
	}

	var i int32
	i, ok = args[1].(int32)
	if !ok {
		t.Error(ok)
		return
	}
	if i != 42 {
		t.Error(i)
		return
	}
}

func TestParseResponse(t *testing.T) {
	source := xml.Header
	source += `<methodResponse>
                   <params><param><value><int>42</int></value>
                       </param>
                   </params>
               </methodResponse>`
	buffer := bytes.NewBufferString(source)
	decoder := xml.NewDecoder(buffer)
	ok, result, e := parseResponse(decoder)
	if e != nil {
		t.Error(e)
		return
	}
	if !ok {
		t.Error(e)
		return
	}
	var i int32
	i, ok = result.(int32)
	if !ok {
		t.Error(e)
		return
	}
	if i != 42 {
		t.Error(e)
		return
	}
}

func TestParseResponse2(t *testing.T) {
	source := `<?xml version="1.0"?>
<methodResponse><params><param>
<value><array><data>
  <value><i4>1</i4></value>
  <value></value>
  <value><array><data>
    <value>TCPROS</value>
    <value>hedgehog</value>
    <value><i4>52060</i4></value>
  </data></array></value>
</data></array></value>
</param></params></methodResponse>`
	buffer := bytes.NewBufferString(source)
	decoder := xml.NewDecoder(buffer)
	ok, result, e := parseResponse(decoder)
	if e != nil {
		t.Error(e)
		return
	}
	if !ok {
		t.Error(e)
		return
	}
	var outer []interface{}
	outer, ok = result.([]interface{})
	if !ok {
		t.Error("Result should be an array.")
	}
	if len(outer) != 3 {
		t.Errorf("Array len was %d, should be 3.", len(outer))
	}
	var i int32
	i, ok = outer[0].(int32)
	if !ok {
		t.Error("First elem should be int.")
	} else if i != 1 {
		t.Errorf("First elem should be 1, was %d.", i)
	}
}

func TestParseFault(t *testing.T) {
	source := xml.Header
	source += `<methodResponse>
                   <fault>
                       <value>
                           <struct>
                               <member>
                                   <name>faultCode</name>
                                   <value><int>42</int></value>
                               </member>
                               <member>
                                   <name>faultString</name>
                                   <value><string>failed</string></value>
                               </member>
                           </struct>
                       </value>
                   </fault>
               </methodResponse>`
	buffer := bytes.NewBufferString(source)
	decoder := xml.NewDecoder(buffer)
	ok, result, e := parseResponse(decoder)
	if e != nil {
		t.Error(e)
		return
	}
	if ok {
		t.Error(e)
		return
	}
	var m map[string]interface{}
	m, ok = result.(map[string]interface{})
	if !ok {
		t.Error(fmt.Sprintf("expected map from string to interface, got: %v with value '%v'",
			reflect.TypeOf(result), result))
		return
	}

	if len(m) != 2 {
		t.Error(m)
		return
	}

	v := m["faultCode"]
	var i int32
	i, ok = v.(int32)
	if !ok {
		t.Error(ok)
		return
	}
	if i != 42 {
		t.Error(i)
		return
	}

	v = m["faultString"]
	var s string
	s, ok = v.(string)
	if !ok {
		t.Error(ok)
		return
	}
	if s != "failed" {
		t.Error(s)
		return
	}
}

func TestClient(t *testing.T) {
	masterURI := os.Getenv("ROS_MASTER_URI")
	t.Log("Master URI: ", masterURI)

	value, e := Call(masterURI, "getPublishedTopics", "not_a_node", "")
	if e != nil {
		t.Error(e)
	}
	t.Log(value)

	value, e = Call(masterURI, "getTopicTypes", "not_a_node")
	if e != nil {
		t.Error(e)
	}
	t.Log(value)

	value, e = Call(masterURI, "getSystemState", "not_a_node")
	if e != nil {
		t.Error(e)
	}
	t.Log(value)

	value, e = Call(masterURI, "getUri", "not_a_node")
	if e != nil {
		t.Error(e)
	}
	t.Log(value)
}

type myDispatcher struct {
	X int32
}

func (h *myDispatcher) addTwoInts(a int32, b int32) (int32, error) {
	c := h.X * (a + b)
	return c, nil
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", ":19937")
	if err != nil {
		panic(err)
		return
	}
	d := myDispatcher{2}
	m := map[string]Method{"addTwoInts": d.addTwoInts}
	handler := NewHandler(m)
	go http.Serve(listener, handler)

	result, e := Call("http://localhost:19937", "addTwoInts", 1, 2)
	if e != nil {
		t.Error(e)
		return
	}
	i, ok := result.(int32)
	if !ok {
		t.Error(ok)
		return
	}
	if i != 6 {
		t.Error(i)
	}

	listener.Close()
	handler.WaitForShutdown()
}
