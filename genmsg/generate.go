package genmsg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"text/template"
)

var msgTemplate = `
//#############################################################################
// Template Context:
// -----------------------------------------------
// file_name_in : String
//     Source file
// spec : msggen.MsgSpec
//     Parsed specification of the .msg file
// md5sum : String
//     MD5Sum of the .msg specification
//#############################################################################
package {{ .Package }}
// Automatically generated from the message definition "{{ .FullName }}"
import (
    "bytes"
{{ if .BinaryRequired }}
    "encoding/binary"
{{ end }}
    "ros"
{{ range .Imports }}
	"{{ . }}"
{{ end }}
)

{{ if gt (len .Constants) 0 }}
const (
{{ range .Constants }}
	{{ if eq .Type "string" }}
    {{ .GoName }} {{ .Type }} = "{{ .Value }}"
	{{ else }}
	{{ .GoName }} {{ .Type }} = {{ .Value }}
	{{ end }}
{{ end }}
)
{{ end }}


type _Msg{{ .ShortName }} struct {
    text string
    name string
    md5sum string
}

func (t *_Msg{{ .ShortName }}) Text() string {
    return t.text
}

func (t *_Msg{{ .ShortName }}) Name() string {
    return t.name
}

func (t *_Msg{{ .ShortName }}) MD5Sum() string {
    return t.md5sum
}

func (t *_Msg{{ .ShortName }}) NewMessage() ros.Message {
    m := new({{ .ShortName }})
{{- range .Fields }}
	name = {{ .GoName }}
{{-     if .IsArray }}
{{-         if eq .ArrayLen -1 }}
	m.{{ .GoName }}} = nil
{{-         else }}
	for i := 0; i < {{ .ArrayLen }}; i++ {
		m.{{ .GoName }}[i]  = {{ .ZeroValue }}
	}
{{-         end}}
{{-     else }}
	m.{{ .GoName }} = {{ .ZeroValue }}
{{-     end }}
{{- end }}
    return m
}

var (
    Msg{{ .ShortName }} = &_Msg{{ .ShortName }} {
        ` + "`" + `{{ .Text }}` + "`" + `,
        "{{ .FullName }}",
        "{{ .MD5Sum }}",
    }
)

type {{ .ShortName }} struct {
{{- range .Fields }}
	{{.GoType }} {{ .GoName }}
{{- end }}
}

func (m *{{ .ShortName }}) Type() ros.MessageType {
	return Msg{{ .ShortName }}
}

func (m *{{ .ShortName }}) Serialize(buf *bytes.Buffer) error {
    var err error = nil
{{- range .Fields }}
{{-     if .IsArray }}
    binary.Write(buf, binary.LittleEndian, uint32(len(m.{{ .GoName }})))
    for _, e := range m.{{ .GoName }} {
{{-         if .IsBuiltin }}
{{-             if eq .Type "string" }}
        binary.Write(buf, binary.LittleEndian, uint32(len([]byte(e))))
        buf.Write([]byte(e))
{{-            else }}
{{-                if or (eq .Type "time") (eq .Type "duration") }}
        binary.Write(buf, binary.LittleEndian, e.Sec)
        binary.Write(buf, binary.LittleEndian, e.NSec)
{{-                else }}
        binary.Write(buf, binary.LittleEndian, e)
{{-                end }}
{{-             end }}
{{-         else }}
        if err = e.Serialize(buf); err != nil {
            return err
        }
{{-         end }}
    }
{{-     else }}
{{-         if .IsBuiltin }}
{{-             if eq .Type "string" }}
    binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.{{ .GoName }}))))
    buf.Write([]byte(m.{{ .GoName }}))
{{-             else }}
{{-                 if or (eq .Type "time") (eq .Type "duration") }}
    binary.Write(buf, binary.LittleEndian, m.{{ .GoName }}.Sec)
    binary.Write(buf, binary.LittleEndian, m.{{ .GoName }}.NSec)
{{-                 else }}
    binary.Write(buf, binary.LittleEndian, m.{{ .GoName }})
{{-                 end }}
{{-             end }}
{{-         else }}
    if err = m.{{ .GoName }}.Serialize(buf); err != nil {
        return err
    }
{{-         end }}
{{-     end }}
{{- end }}
    return err
}


func (m *{{ .ShortName }}) Deserialize(buf *bytes.Reader) error {
    var err error = nil
{{- range .Fields }}
{{-    if .IsArray }}
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
{{-        if lt .ArrayLen 0 }}
        m.{{ .GoName }} = make([]{{ .GoType }}), int(size))
{{-        end }}
        for i := 0; i < int(size); i++ {
{{-          if .IsBuiltin }}
{{-              if eq .Type "string" }}
            {
                var size uint32
                if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
                    return err
                }
                data := make([]byte, int(size))
                if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
                    return err
                }
                m.{{ .GoName }})[i] = string(data)
            }
{{-              else }}
{{- 					if or (eq .Type "time") (eq .Type "duration") }}
            {
                if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}[i].Sec); err != nil {
                    return err
                }

                if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}[i].NSec); err != nil {
                    return err
                }
            }
{{-                  else }}
            if err = binary.Read(buf, binary.LittleEndian, &m.{{ .Name }}[i]); err != nil {
                return err
            }
{{-                  end }}
{{-              end }}
{{-          else }}
            if err = m.{{ .GoName }}[i].Deserialize(buf); err != nil {
                return err
            }
{{-      	end }}
        }
    }
{{-    else }}
{{-        if .IsBuiltin }}
{{-            if eq .Type "string" }}
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        data := make([]byte, int(size))
        if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
            return err
        }
        m.{{ .GoName }} = string(data)
    }
{{-            else }}
{{-            		if or (eq .Type "time") (eq .Type "duration") }}
    {
        if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}.Sec); err != nil {
            return err
        }

        if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}.NSec); err != nil {
            return err
        }
    }
{{-            		else }}
    if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}); err != nil {
        return err
    }
{{-         			end }}
{{-            end }}
{{-        else }}
    if err = m.{{ .GoName }}.Deserialize(buf); err != nil {
        return err
    }
{{-    	  end }}
{{-    end }}
{{- end }}
    return err
}
`

var srvTemplate = `
//##############################################################################
//# Template Context:
//# -----------------------------------------------
//# file_name_in : String
//#     Source file
//# spec : msggen.MsgSpec
//#     Parsed specification of the .msg file
//# md5sum : String
//#     MD5Sum of the .msg specification
//##############################################################################
package {{ .Pacakge }}
// Automatically generated from {{ .Filename }}
import (
    "ros"
)

// Service type metadata
type _Srv{{ .ShortName }} struct {
    name string
    md5sum string
    text string
    reqType ros.MessageType
    resType ros.MessageType
}

func (t *_Srv{{ .ShortName }}) Name() string { return t.name }
func (t *_Srv{{ .ShortName }}) MD5Sum() string { return t.md5sum }
func (t *_Srv{{ .ShortName }}) Text() string { return t.text }
func (t *_Srv{{ .ShortName }}) RequestType() ros.MessageType { return t.reqType }
func (t *_Srv{{ .ShortName }}) ResponseType() ros.MessageType { return t.resType }
func (t *_Srv{{ .ShortName }}) NewService() ros.Service {
    return new({{ .ShortName }})
}

var (
    Srv{{ .ShortName }} = &_Srv{{ .ShortName }} {
        "{{ .FullName }}",
        "{{ .MD5Sum }}",
        ` + "`" + `{{ .Text }}` + "`" + `,
        Msg{{ .ShortName }}Request,
        Msg{{ .ShortName }}Response,
    }
)


type {{ .ShortName }} struct {
    Request {{ .ShortName }}Request
    Response {{ .ShortName }}Response
}

func (s *{{ .ShortName }}) NewRequest() ros.Message { return &s.Request }
func (s *{{ .ShortName }}) NewResponse() ros.Message { return &s.Response }
`

type GenMsgContext struct {
	MsgSpec
	MD5Sum         string
	BinaryRequired bool
	Imports        []string
}

func GenerateMessage(spec *MsgSpec) (string, error) {
	var context GenMsgContext
	context.Fields = spec.Fields
	context.Constants = spec.Constants
	context.Text = spec.Text
	context.FullName = spec.FullName
	context.ShortName = spec.ShortName
	context.Package = spec.Package

	hash := md5.New()
	context.MD5Sum = hex.EncodeToString(hash.Sum([]byte(context.Text)))

	for _, field := range context.Fields {
		if isPrimitiveType(field.Type) {
			context.BinaryRequired = true
			break
		}
	}

	tmpl, err := template.New("msg").Parse(msgTemplate)

	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, context)

	if err != nil {
		panic(err)
	}
	return buffer.String(), err
}

// func GenerateService(spec *SrvSpec) string, error {
// 	tmpl, err := template.New("msg").Parse(srcTemplate)
//
// 	return nil, nil
// }
