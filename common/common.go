package common

import (
	"bufio"
	"strings"
)

const (
	MESSAGE = "message"
	ONEOF = "oneof"
	ENUM = "enum"
	SPACE = " "
	SPACE_RUNE = ' '
	REPEATED = "repeated"
	OPTIONAL = "optional"
	RESERVED = "reserved"
	MAP = "map"
)

var KNOWN_TYPES = map[string]string{
	"bool": "boolean",
	"string": "string",
	"float": "number",
	"double": "number",
	"int32": "number",
	"uint32": "number",
	"bytes": "Buffer",
	"int64": "Int64",
	"uint64": "Int64",
	"google.protobuf.Int64Value": "Int64Value",
	"google.protobuf.FloatValue": "FloatValue",
	"google.protobuf.DoubleValue": "DoubleValue",
	"google.protobuf.BoolValue": "BoolValue",
	"google.protobuf.StringValue": "StringValue",
	"google.protobuf.Timestamp": "Timestamp",
	"google.protobuf.Duration": "Duration",
	"google.type.TimeOfDay": "TimeOfDay",
	"google.type.DayOfWeek": "DayOfWeek",
	"google.protobuf.FieldMask": "FieldMask",
	"google.protobuf.Empty": "Empty",
	"google.protobuf.Any": "Any",
	"google.protobuf.Struct": "Struct",
	"google.protobuf.Value": "Value",
	"google.protobuf.ListValue": "ListValue",
}


func nextNotSpaceIdx (spaceIdx int, lineText string) int {
	i := spaceIdx + 1;
	for ; i < len(lineText) && lineText[i] == SPACE_RUNE; i++ {}
	return i
}

func nextSpaceIdx (startIdx int, lineText string) int {
	i := startIdx + 1
	for ; i < len(lineText) && lineText[i] != SPACE_RUNE; i++ {}
	return i
}

func GetSecondWord (lineText string, firstWord string) string {
	startIdx := strings.IndexAny(lineText, firstWord) + len(firstWord)
	startIdx = nextNotSpaceIdx(startIdx, lineText)

	endIdx := nextSpaceIdx(startIdx, lineText)

	return lineText[startIdx: endIdx]
}

type MyScanner struct {
    *bufio.Scanner
}

func (s *MyScanner) TrimedText() string {
	return strings.TrimSpace(s.Text())
}
