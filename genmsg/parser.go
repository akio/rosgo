package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	Sep    = "/"
	MsgDir = "msg"
	SrvDir = "srv"
	ExtMsg = ".msg"
	ExtSrv = ".msg"

	ConstChar   = "="
	CommentChar = "#"
	IoDelim     = "---"
)

type SyntaxError struct {
	FullName string
	Line     int
	Message  string
}

func NewSyntaxError(fullName string, line int, message string) *SyntaxError {
	self := &SyntaxError{}
	self.FullName = fullName
	self.Line = line
	self.Message = message
	return self
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("[%s@%d] %s", e.FullName, e.Line, e.Message)
}

/// Convert constant literal to a Go object
/// Original implementation (genmsg) depends on Python's literal syntax.
func convertConstantValue(fieldType string, valueLiteral string) (interface{}, error) {
	switch fieldType {
	case "float32":
		result, e := strconv.ParseFloat(valueLiteral, 32)
		return float32(result), e
	case "float64":
		return strconv.ParseFloat(valueLiteral, 64)
	case "string":
		return strings.TrimSpace(valueLiteral), nil
	case "byte":
		result, e := strconv.ParseInt(valueLiteral, 0, 8)
		return int8(result), e
	case "int8":
		result, e := strconv.ParseInt(valueLiteral, 0, 8)
		return int8(result), e
	case "int16":
		result, e := strconv.ParseInt(valueLiteral, 0, 16)
		return int16(result), e
	case "int32":
		result, e := strconv.ParseInt(valueLiteral, 0, 32)
		return int32(result), e
	case "int64":
		return strconv.ParseInt(valueLiteral, 0, 64)
	case "char":
		result, e := strconv.ParseUint(valueLiteral, 0, 8)
		return uint8(result), e
	case "uint8":
		result, e := strconv.ParseUint(valueLiteral, 0, 8)
		return uint8(result), e
	case "uint16":
		result, e := strconv.ParseUint(valueLiteral, 0, 16)
		return uint16(result), e
	case "uint32":
		result, e := strconv.ParseUint(valueLiteral, 0, 32)
		return uint32(result), e
	case "uint64":
		return strconv.ParseUint(valueLiteral, 0, 64)
	case "bool":
		// The spec of ROS message doesn't specify boolean literal exactly.
		// genmsg implementation determines true/false based Python's eval() and accepts any valid Python expression.
		if valueLiteral == "None" || valueLiteral == "False" {
			return false, nil
		} else if valueLiteral == "True" {
			return true, nil
		} else if val, e := strconv.ParseUint(valueLiteral, 10, 0); e == nil {
			return val != 0, nil
		} else {
			return nil, fmt.Errorf("Inavalid constant literal for bool: [%s]", valueLiteral)
		}
	default:
		return nil, fmt.Errorf("Invalid constant type: [%s]", fieldType)
	}
}

func packageResourceName(name string) (string, string, error) {
	const Separator = "/"
	if strings.Contains(name, Separator) {
		components := strings.Split(name, Separator)
		if len(components) == 2 {
			return components[0], components[1], nil
		} else {
			return "", "", fmt.Errorf("Invalid name %s", name)
		}
	} else {
		return "", name, nil
	}
}

func stripComment(line string) string {
	return strings.TrimSpace(strings.Split(line, CommentChar)[0])
}

func loadConstantLine(line string) (*Constant, error) {
	cleanLine := stripComment(line)
	sepIndex := strings.IndexFunc(cleanLine, unicode.IsSpace)
	if sepIndex < 0 {
		return nil, fmt.Errorf("Could not find a constant name after the type name")
	}

	fieldType := cleanLine[:sepIndex]
	if !isValidConsantType(fieldType) {
		return nil, fmt.Errorf("[%s] is not a legal constant type", fieldType)
	}

	var name, valueText string
	if fieldType == "string" {
		// Strings contain anything to the right of the equal sign, no comments allowd.
		sepIndex := strings.IndexFunc(line, unicode.IsSpace)
		if sepIndex < 0 {
			return nil, fmt.Errorf("Could not find a constant name after the type name")
		}
		keyValue := line[sepIndex:]
		kvSplits := strings.SplitN(keyValue, "=", 2)
		if len(kvSplits) != 2 {
			return nil, fmt.Errorf("A constant definition requires its value")
		}
		name = strings.TrimSpace(kvSplits[0])
		valueText = strings.TrimLeftFunc(kvSplits[1], unicode.IsSpace)
	} else {
		keyValue := strings.TrimSpace(cleanLine[sepIndex:])
		kvSplits := strings.SplitN(keyValue, "=", 2)
		if len(kvSplits) != 2 {
			return nil, fmt.Errorf("A constant definition requires its value")
		}
		name = strings.TrimSpace(kvSplits[0])
		valueText = strings.TrimSpace(kvSplits[1])
	}

	value, e := convertConstantValue(fieldType, valueText)
	if e != nil {
		return nil, e
	}
	return NewConstant(fieldType, name, value, valueText), nil
}

func loadFieldLine(line string, packageName string) (*Field, error) {
	cleanLine := stripComment(line)
	lineSplits := strings.SplitN(cleanLine, " ", 2)
	if len(lineSplits) != 2 {
		return nil, fmt.Errorf("Invalid declaration: %s", line)
	}
	fieldType := strings.TrimSpace(lineSplits[0])
	name := strings.TrimSpace(lineSplits[1])
	if !isValidMsgFieldName(name) {
		return nil, fmt.Errorf("%s is not a legal message field name", name)
	}
	if !isValidMsgType(fieldType) {
		return nil, fmt.Errorf("%s is not a legal message field type", fieldType)
	}
	if len(packageName) > 0 && !strings.Contains(fieldType, Sep) {
		if fieldType == HeaderType {
			fieldType = HeaderFullName
		} else if !isBuiltinType(baseMsgType(fieldType)) {
			fieldType = fmt.Sprintf("%s/%s", packageName, fieldType)
		}
	} else if fieldType == HeaderType {
		fieldType = HeaderFullName
	}
	pkg, baseType, isArray, arrayLen, err := parseType(fieldType)
	if err != nil {
		return nil, err
	}

	return NewField(pkg, baseType, name, isArray, arrayLen), nil
}
