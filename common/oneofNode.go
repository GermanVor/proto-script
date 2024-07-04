package common

import (
	"fmt"
	"strings"
)

type OneofNode struct {
	OneofName string
	Body []*VariableNode
}

func InitOneofNode(oneofName string) *OneofNode {
	return &OneofNode{
		OneofName: oneofName,
		Body: []*VariableNode{},
	}
}

func (oNode *OneofNode) AddLine(trimedLineText string) {
	// new line case
	if trimedLineText == "" {
		return
	}

	// comment case
	if strings.Index(trimedLineText, "//") == 0 {
		return
	}

	// end case
	if strings.Contains(trimedLineText, "}") {
		return
	}

	vNode := InitVariableNode(trimedLineText)

	oNode.Body = append(
		oNode.Body,
		vNode,
	)
}

func (oNode *OneofNode) String() string {
	res := fmt.Sprintf("export type %s = OneOf<'%s', {\n", oNode.OneofName, oNode.OneofName)

	for _, vNode := range oNode.Body {
		res += vNode.String()
		res += "\n"
	}

	res += "}>;"

	return res
}
