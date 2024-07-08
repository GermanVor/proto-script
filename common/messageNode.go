package common

import (
	"fmt"
	"slices"
	"strings"
)

type NamespaceBodyType interface {
	String() string
}

type MessageNode struct {
	MessageName string
	MessageType string
	body []*VariableNode
	NamespaceTypeList []string
	NamespaceBody []NamespaceBodyType
	substitutionMap SubstitutionMap
}

func InitMessageNode(messageName string, substitutionMap SubstitutionMap) *MessageNode {
	return &MessageNode{
		MessageName: messageName,
		MessageType: messageName,
		body: []*VariableNode{},
		NamespaceTypeList: []string{},
		NamespaceBody: []NamespaceBodyType{},
		substitutionMap: substitutionMap,
	}
}

func (mNode *MessageNode) AddOneofNode(oNode *OneofNode) {
	mNode.NamespaceTypeList = append(mNode.NamespaceTypeList, oNode.OneofName)
	mNode.NamespaceBody = append(mNode.NamespaceBody, oNode)
}

func (mNode *MessageNode) AddEnumNode(eNode *EnumNode) {
	mNode.NamespaceTypeList = append(mNode.NamespaceTypeList, eNode.EnumName)
	mNode.NamespaceBody = append(mNode.NamespaceBody, eNode)
}

func (mNode *MessageNode) AddMessageNode(mNodeToAdd *MessageNode) {
	mNode.NamespaceTypeList = append(mNode.NamespaceTypeList, mNodeToAdd.MessageName)
	mNode.NamespaceBody = append(mNode.NamespaceBody, mNodeToAdd)
}

func (mNode *MessageNode) AddLine(trimedLineText string) {
	// end case
	if strings.Contains(trimedLineText, "}") {
		return
	}

	if vNode := InitVariableNode(trimedLineText); vNode != nil {
		mNode.body = append(mNode.body, vNode)
	}
}

func (mNode *MessageNode) String() string {
	res := ""

	oneofName := ""

	isNamespace := len(mNode.NamespaceBody) > 0

	if isNamespace {
		res += "// eslint-disable-next-line @typescript-eslint/no-namespace\n"
		res += fmt.Sprintf("export namespace %s {\n", mNode.MessageName)

		for _, line := range mNode.NamespaceBody {
			res += line.String()
			res += "\n"

			if oneofNode, ok := line.(*OneofNode); ok {
				oneofName = oneofNode.OneofName
			}
		}

		res += "};\n"
	}

	if isNamespace {
		res += "// eslint-disable-next-line @typescript-eslint/no-redeclare\n"
	}

	if oneofName == "" {
		res += fmt.Sprintf("export type %s = {\n", mNode.MessageName)
	} else {
		fullOneofTypeName := fmt.Sprintf("%s.%s", mNode.MessageName, oneofName)
		res += fmt.Sprintf("export type %s = %s & {\n", mNode.MessageName, fullOneofTypeName)
	}

	for _, vNode := range mNode.body {
		res += vNode.String()
		res += "\n"
	}

	res += "};"

	return res
}

func (mNode *MessageNode) PrepareVars(importSourceMap ImportMap) {
	getSubstitutionType := func (vType string) (string, bool) {
		if mNode.substitutionMap == nil {
			return vType, false
		}

		if tsValue, ok := mNode.substitutionMap[vType]; ok {
			if tsValue.ImportSource != "" {
				if importSourceMap[tsValue.ImportSource] == nil {
					importSourceMap[tsValue.ImportSource] = make(Set)
				}

				importSourceMap[tsValue.ImportSource][tsValue.Substitution] = true
			}

			return tsValue.Substitution, true
		}

		return vType, false
	}

	getRealType := func (vType string) (string, bool) {
		if slices.Contains(mNode.NamespaceTypeList, vType) {
			vType = fmt.Sprintf("%s.%s", mNode.MessageName, vType)
			return vType, true
		}

		return getSubstitutionType(vType)
	}

	for _, vNode := range mNode.body {
		if vNode.Type.stringType != nil {
			stringType, replaced := getRealType(*vNode.Type.stringType)

			if replaced {
				vNode.Type = VariableNodeType{
					stringType: &stringType,
				}
			}

			continue
		}

		if vNode.Type.mapType != nil {
			mapKeyType, mapKeyTypeReplaced := getRealType(vNode.Type.mapType[0])
			mapValueType, mapValueTypeReplaced := getRealType(vNode.Type.mapType[1])

			if mapKeyTypeReplaced || mapValueTypeReplaced {
				mapType := [2]string{mapKeyType, mapValueType}
				vNode.Type = VariableNodeType{
					mapType: &mapType,
				}
			}

			continue
		}

		fmt.Println(vNode)
	}
}
