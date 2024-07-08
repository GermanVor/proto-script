package common

import (
	"fmt"
	"strings"
)

type EnumNode struct {
	EnumName string
	Body []string
}

func InitEnumNode(enumName string) *EnumNode {
	return &EnumNode{
		EnumName: enumName,
		Body: []string{},
	}
}

func (eNode *EnumNode) AddLine(trimedLineText string) {
	// end case
	if strings.Contains(trimedLineText, "}") {
		return
	}

	nodeValue := trimedLineText[0: strings.Index(trimedLineText, SPACE)]
	eNode.Body = append(eNode.Body, fmt.Sprintf("%s = '%s',", nodeValue, nodeValue))
}

func (eNode *EnumNode) String() string {
	res := fmt.Sprintf("export enum %s {\n", eNode.EnumName)

	for _, line := range eNode.Body {
		res += line
		res += "\n"
	}

	res += "}"

	return res
}
