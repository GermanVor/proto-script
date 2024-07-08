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
	// end case
	if strings.Contains(trimedLineText, "}") {
		return
	}

	// option
	if strings.Index(trimedLineText, OPTION) == 0 {
		return
	}

	vNode := InitVariableNode(trimedLineText)

	oNode.Body = append(
		oNode.Body,
		vNode,
	)
}

func (oNode *OneofNode) String() string {
	res := fmt.Sprintf(
		"export type %s = %s<'%s', {\n",
		oNode.OneofName,
		ONEOF_GENERIC_NAME,
		oNode.OneofName,
	)

	for _, vNode := range oNode.Body {
		res += vNode.String()
		res += "\n"
	}

	res += "}>;"

	return res
}
