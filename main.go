package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GermanVor/proto-script/common"
)

// true - если внутри коллбека снова вызывается forNode
// false - если внутри коллбека не вызывается forNode
type ForNodeCallBack func(line string, i int) bool

func forNode (s *common.MyScanner, callBack ForNodeCallBack) int {
	trimedLineText := s.TrimedText()
	openBracketCount := 0;

	i := 0

	for {
		isTrigeredForNodeInside := callBack(trimedLineText, i)

		if strings.Contains(trimedLineText, "{") && !isTrigeredForNodeInside {
			openBracketCount += 1
		}

		if strings.Contains(trimedLineText, "}") {
			openBracketCount -= 1
		}

		if openBracketCount <= 0 || !s.Scan() {
			break
		}

		trimedLineText = s.TrimedText()
		i++
	}

	return openBracketCount
}

func parseOneof (s *common.MyScanner) *common.OneofNode {
	oneofName := common.GetSecondWord(s.TrimedText(), common.ONEOF)
	oNode := common.InitOneofNode(oneofName)

	forNode(s, func(trimedLineText string, i int) bool {
		if i == 0 {
			return false
		}

		oNode.AddLine(trimedLineText)

		return false
	})

	return oNode
}

func parseEnum (s *common.MyScanner) *common.EnumNode {
	enumName := common.GetSecondWord(s.TrimedText(), common.ENUM)
	eNode := common.InitEnumNode(enumName)

	forNode(s, func(trimedLineText string, i int) bool {
		if i == 0 {
			return false
		}

		eNode.AddLine(trimedLineText)
		return false
	})

	return eNode
}

func parseMessage (
	s *common.MyScanner,
	substitutionMap common.SubstitutionMap,
	importSourceMap common.ImportMap,
) *common.MessageNode {
	messageName := common.GetSecondWord(s.TrimedText(), common.MESSAGE)
	mNode := common.InitMessageNode(messageName, substitutionMap)

	forNode(s, func(trimedLineText string, i int) bool {
		if i == 0 {
			return false
		}

		if strings.Index(trimedLineText, common.ONEOF) == 0 {
			if substitution, ok := substitutionMap["oneof"]; ok {
				if substitution.ImportSource != "" {
					importSourceMap[substitution.ImportSource] = make(map[string]bool)
					importSourceMap[substitution.ImportSource][common.ONEOF_GENERIC_NAME] = true
				}
			}

			oNode := parseOneof(s)
			mNode.AddOneofNode(oNode)
			return true
		}

		if strings.Index(trimedLineText, common.ENUM) == 0 {
			eNode := parseEnum(s)
			mNode.AddEnumNode(eNode)
			return true
		}

		if strings.Index(trimedLineText, common.MESSAGE) == 0 {
			mNodeToAdd := parseMessage(s, substitutionMap, importSourceMap)
			mNode.AddMessageNode(mNodeToAdd)
			return true
		}

		mNode.AddLine(trimedLineText)
		return false
	})

	mNode.PrepareVars(importSourceMap)

	return mNode
}

func main() {
	var protoPath string
	var key string

	flag.StringVar(&protoPath, "p", "", "Absolute proto path")
	flag.StringVar(&key, "k", "datasphere", "Config subtitution key")

	flag.Parse()

	if protoPath == "" {
		log.Fatal("Specify proto path !")
	}

	substitutionMap, _ := common.ParseConfig(key)

	file, err := os.Open(protoPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	s := &common.MyScanner{Scanner: bufio.NewScanner(file)}
	mList := make([]*common.MessageNode, 0)

	importSourceMap := common.ImportMap{};

	for s.Scan() {
		lineText := s.TrimedText()

		if strings.Contains(lineText, common.MESSAGE) {
			mNode := parseMessage(
				s,
				substitutionMap,
				importSourceMap,
			)
			mList = append(mList, mNode)
		}
	}

	fmt.Println(common.ImportMapToString(importSourceMap))

	for _, mNode := range mList {
		fmt.Println(mNode.String())
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}
