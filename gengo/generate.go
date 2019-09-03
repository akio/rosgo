package main

import (
	"bytes"
	"text/template"
)

var msgTemplate = `
// Automatically generated from the message definition "{{ .FullName }}.msg"
package {{ .Package }}
import (
    "bytes"
{{- if .BinaryRequired }}
    "encoding/binary"
{{- end }}
    "github.com/edwinhayes/rosgo/ros"
{{- range .Imports }}
	"{{ . }}"
{{- end }}
)

{{- if gt (len .Constants) 0 }}
const (
{{- range .Constants }}
	{{- if eq .Type "string" }}
    {{ .GoName }} {{ .Type }} = "{{ .Value }}"
	{{- else }}
	{{ .GoName }} {{ .Type }} = {{ .Value }}
	{{- end }}
{{- end }}
)
{{- end }}


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
{{-     if .IsArray }}
{{-         if eq .ArrayLen -1 }}
	m.{{ .GoName }} = []{{ .GoType }}{}
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
{{-     if .IsArray }}
{{-         if eq .ArrayLen -1 }}
	{{ .GoName }} []{{ .GoType }}` + " `rosmsg:\"{{ .Name }}:{{ .Type }}[]\"`" + `
{{-         else }}
	{{ .GoName }} [{{ .ArrayLen }}]{{ .GoType }}` + " `rosmsg:\"{{ .Name }}:{{ .Type }}[{{ .ArrayLen }}]\"`" + `
{{-         end }}
{{-     else }}
	{{ .GoName }} {{ .GoType }}` + " `rosmsg:\"{{ .Name }}:{{ .Type }}\"`" + `
{{-     end }}
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
        m.{{ .GoName }} = make([]{{ .GoType }}, int(size))
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
            if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}[i]); err != nil {
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
// Automatically generated from the message definition "{{ .FullName }}.srv"
package {{ .Package }}
import (
    "github.com/akio/rosgo/ros"
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

func (s *{{ .ShortName }}) ReqMessage() ros.Message { return &s.Request }
func (s *{{ .ShortName }}) ResMessage() ros.Message { return &s.Response }
`

type MsgGen struct {
	MsgSpec
	BinaryRequired bool
	Imports        []string
}

func (gen *MsgGen) analyzeImports() {
	for _, field := range gen.Fields {
		if len(field.Package) == 0 {
			gen.BinaryRequired = true
		} else {
			found := false
			for _, imp := range gen.Imports {
				if imp == field.Package {
					found = true
					break
				}
			}
			if !found {
				gen.Imports = append(gen.Imports, field.Package)
			}
		}
	}
}

func GenerateMessage(context *MsgContext, spec *MsgSpec) (string, error) {
	var gen MsgGen
	gen.Fields = spec.Fields
	gen.Constants = spec.Constants
	gen.Text = spec.Text
	gen.FullName = spec.FullName
	gen.ShortName = spec.ShortName
	gen.Package = spec.Package
	gen.MD5Sum = spec.MD5Sum

	gen.analyzeImports()

	tmpl, err := template.New("msg").Parse(msgTemplate)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, gen)
	if err != nil {
		return "", err
	}
	return buffer.String(), err
}

func GenerateService(context *MsgContext, spec *SrvSpec) (string, string, string, error) {
	reqCode, err := GenerateMessage(context, spec.Request)
	if err != nil {
		return "", "", "", err
	}
	resCode, err := GenerateMessage(context, spec.Response)
	if err != nil {
		return "", "", "", err
	}

	tmpl, err := template.New("srv").Parse(srvTemplate)
	if err != nil {
		return "", "", "", err
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, spec)
	if err != nil {
		return "", "", "", err
	}
	return buffer.String(), reqCode, resCode, err
}
