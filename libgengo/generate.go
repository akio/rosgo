package libgengo

import (
	"bytes"
	"text/template"
)

var import_path *string

var msgTemplate = `
// Package {{ .Package }} is automatically generated from the message definition "{{ .FullName }}.msg"
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

    "github.com/edwinhayes/rosgo/ros"
)

{{- if gt (len .Constants) 0 }}
const (
{{- range .Constants }}
	{{- if eq .Type "string" }}
    {{ $.ShortName }}_{{ .GoName }} {{ .Type }} = "{{ .Value }}"
	{{- else }}
	{{ $.ShortName }}_{{ .GoName }} {{ .Type }} = {{ .Value }}
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

func (m *{{ .ShortName }}) GetType() ros.MessageType {
	return Msg{{ .ShortName }}
}

func (m *{{ .ShortName }}) Serialize(buf *bytes.Buffer) error {
    var err error
{{- range .Fields }}
{{-     if .IsArray }}
{{-        if lt .ArrayLen 0 }}
    binary.Write(buf, binary.LittleEndian, uint32(len(m.{{ .GoName }})))
{{-        end }}
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

{{-        if lt .ArrayLen 0 }}
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        m.{{ .GoName }} = make([]{{ .GoType }}, int(size))
        for i := 0; i < int(size); i++ {
{{-        else }}
        for i :=0; i < {{ .ArrayLen }}; i++ {
{{-        end }}
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
                m.{{ .GoName }}[i] = string(data)
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

{{- if .IsAction }}
{{- range .Fields }}
{{-     if or (eq .GoName "Goal") (eq .GoName "Feedback") (eq .GoName "Result") }} 
            
            func (m *{{ $.ShortName }}) Get{{ .GoName }}() ros.Message {
                return &m.{{ .GoName }}
            }

            func (m *{{ $.ShortName }}) Set{{ .GoName }}(s ros.Message) {
                msg := s.(*{{ .GoType }})
                m.{{ .GoName }} = *msg
            }

{{-     else}}
            
            func (m *{{ $.ShortName }}) Get{{ .GoName }}() {{ .GoType }} {
                return m.{{ .GoName }}
            }

            func (m *{{ $.ShortName }}) Set{{ .GoName }}(s {{ .GoType }}) {
                m.{{ .GoName }} = s
            }

{{-     end}}
{{- end }}
{{- end}}
`

var srvTemplate = `
// Package {{ .Package }} is automatically generated from the message definition "{{ .FullName }}.srv"
package {{ .Package }}
import (
    "github.com/fetchrobotics/rosgo/ros"
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

var actionTemplate = `
// Automatically generated from the message definition "{{ .FullName }}.action"
package {{ .Package }}
import (
    "github.com/fetchrobotics/rosgo/actionlib"
    "github.com/fetchrobotics/rosgo/ros"
)

// Service type metadata
type _Action{{ .ShortName }} struct {
    name string
    md5sum string
    text string
    goalType ros.MessageType
    feedbackType ros.MessageType
    resultType ros.MessageType
}

func (t *_Action{{ .ShortName }}) Name() string { return t.name }
func (t *_Action{{ .ShortName }}) MD5Sum() string { return t.md5sum }
func (t *_Action{{ .ShortName }}) Text() string { return t.text }
func (t *_Action{{ .ShortName }}) GoalType() ros.MessageType { return t.goalType }
func (t *_Action{{ .ShortName }}) FeedbackType() ros.MessageType { return t.feedbackType }
func (t *_Action{{ .ShortName }}) ResultType() ros.MessageType { return t.resultType }
func (t *_Action{{ .ShortName }}) NewAction() actionlib.Action {
    return new({{ .ShortName }})
}

var (
    Action{{ .ShortName }} = &_Action{{ .ShortName }} {
        "{{ .FullName }}",
        "{{ .MD5Sum }}",
        ` + "`" + `{{ .Text }}` + "`" + `,
        Msg{{ .ShortName }}ActionGoal,
        Msg{{ .ShortName }}ActionFeedback,
        Msg{{ .ShortName }}ActionResult,
    }
)


type {{ .ShortName }} struct {
    Goal {{ .ShortName }}ActionGoal
    Feedback {{ .ShortName }}ActionFeedback
    Result {{ .ShortName }}ActionResult
}

func (s *{{ .ShortName }}) GetActionGoal() actionlib.ActionGoal         { return &s.Goal }
func (s *{{ .ShortName }}) GetActionFeedback() actionlib.ActionFeedback { return &s.Feedback }
func (s *{{ .ShortName }}) GetActionResult() actionlib.ActionResult     { return &s.Result }
`

type MsgGen struct {
	MsgSpec
	BinaryRequired bool
	IsAction       bool
	Imports        []string
}

func (gen *MsgGen) analyzeImports() {
	fullpath := ""
	if len(*importPath) != 0 {
		fullpath = *importPath + "/"
	}

LOOP:
	for i, field := range gen.Fields {
		if len(field.Package) == 0 {
			gen.BinaryRequired = true
		} else if gen.Package == field.Package {
			gen.Fields[i].GoType = field.Type
			gen.Fields[i].ZeroValue = field.Type + "{}"
		} else {
			for _, imp := range gen.Imports {
				if imp == fullpath+field.Package {
					continue LOOP
				}
			}
			gen.Imports = append(gen.Imports, fullpath+field.Package)
		}

		// Binary is required to read the size of array
		if field.IsArray {
			gen.BinaryRequired = true
		}
	}
}

func GenerateMessage(context *MsgContext, spec *MsgSpec, isAction bool) (string, error) {
	var gen MsgGen
	gen.IsAction = isAction
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
	reqCode, err := GenerateMessage(context, spec.Request, false)
	if err != nil {
		return "", "", "", err
	}
	resCode, err := GenerateMessage(context, spec.Response, false)
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

type ActionCode struct {
	goalCode string
}

func GenerateAction(context *MsgContext, spec *ActionSpec) (actionCode string, codeMap map[string]string, err error) {
	codeMap = make(map[string]string)
	codeMap[spec.Goal.FullName], err = GenerateMessage(context, spec.Goal, false)
	if err != nil {
		return
	}

	codeMap[spec.ActionGoal.FullName], err = GenerateMessage(context, spec.ActionGoal, true)
	if err != nil {
		return
	}
	codeMap[spec.Result.FullName], err = GenerateMessage(context, spec.Result, false)
	if err != nil {
		return
	}
	codeMap[spec.ActionResult.FullName], err = GenerateMessage(context, spec.ActionResult, true)
	if err != nil {
		return
	}
	codeMap[spec.Feedback.FullName], err = GenerateMessage(context, spec.Feedback, false)
	if err != nil {
		return
	}
	codeMap[spec.ActionFeedback.FullName], err = GenerateMessage(context, spec.ActionFeedback, true)
	if err != nil {
		return
	}

	tmpl, err := template.New("action").Parse(actionTemplate)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, spec)
	if err != nil {
		return
	}
	actionCode = buffer.String()

	return
}
