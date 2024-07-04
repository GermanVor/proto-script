package common

import (
	"fmt"
	"strings"
	"unicode"
)

func toCamelCase(s string) string {
	res := ""

	for i, symbol := range s {
		if symbol == '_' {
			continue
		}

		if i > 0 && s[i-1] == '_' {
			res += string(unicode.ToUpper(symbol))
			continue
		}

		res += string(symbol)
	}

	return res
}

func min(value_0, value_1 int) int {
	if value_0 < value_1 {
		return value_0
	}

	return value_1
}

type MapType = [2]string
type VariableNodeType struct{
	stringType *string
	mapType *MapType
}

type VariableNode struct {
	Name string
	Type VariableNodeType
	IsArray bool
	IsOptional bool
}
func InitVariableNode (trimedLineText string) *VariableNode {
	// https://protobuf.dev/programming-guides/proto3/#field-labels

	// optional case
	isOptional := false
	if strings.Index(trimedLineText, OPTIONAL) == 0 {
		startIdx := nextNotSpaceIdx(len(OPTIONAL), trimedLineText)
		trimedLineText = trimedLineText[startIdx:]

		isOptional = true
	}

	// repeated case
	isArray := false
	if strings.Index(trimedLineText, REPEATED) == 0 {
		startIdx := nextNotSpaceIdx(len(REPEATED), trimedLineText)
		trimedLineText = trimedLineText[startIdx:]

		isArray = true
	}

	typeValue := VariableNodeType{}

	// map case
	if strings.Index(trimedLineText, MAP) == 0 {
		// map<string, string> locales = 1;

		keyStartIdx := strings.Index(trimedLineText, "<") + 1
		keyEndIdx := min(
			strings.Index(trimedLineText, SPACE),
			strings.Index(trimedLineText, ","),
		)

		key := trimedLineText[keyStartIdx: keyEndIdx]

		valueStartIdx := nextNotSpaceIdx(keyEndIdx, trimedLineText)
		valueEndIdx := strings.Index(trimedLineText, ">")

		value := trimedLineText[valueStartIdx: valueEndIdx]

		mapType := [2]string{key, value}
		typeValue.mapType = &mapType

		trimedLineText = trimedLineText[nextNotSpaceIdx(valueEndIdx, trimedLineText): ]
	} else {
		// string locales = 1;

		endIdx := strings.Index(trimedLineText, SPACE)
		stringType := trimedLineText[0: endIdx]
		typeValue.stringType = &stringType

		trimedLineText = trimedLineText[nextNotSpaceIdx(endIdx, trimedLineText): ]
	}

	// locales = 1;
	varName := trimedLineText[0: strings.Index(trimedLineText, SPACE)]

	return &VariableNode{
		Name: varName,
		Type: typeValue,
		IsArray: isArray,
		IsOptional: isOptional,
	}
}

func (vNode *VariableNode) String() string {
	res := toCamelCase(vNode.Name)

	if vNode.IsOptional {
		res += "?"
	}

	if vNode.Type.stringType != nil {
		res += fmt.Sprintf(": %s", *vNode.Type.stringType)
	}

	if vNode.Type.mapType != nil {
		res += fmt.Sprintf(
			": Record<%s, %s",
			vNode.Type.mapType[0],
			vNode.Type.mapType[1],
		)

		if vNode.Type.mapType[0] == "string" {
			res += " | undefined"
		}

		res += ">"
	}

	if vNode.IsArray {
		res += "[]"
	}

	res += ";"

	return res
}
