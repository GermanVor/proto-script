package common

import (
	"bufio"
	"fmt"
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
	OPTION = "option"
	RESERVED = "reserved"
	MAP = "map"
	ONEOF_GENERIC_NAME = "OneOf"
)

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

type Set = map[string]bool
type ImportSource = string
type ImportMap = map[ImportSource]Set

func ImportMapToString (importMap ImportMap) string {
	res := ""

	for importSource, typesToImport := range importMap {
		res += "import {\n"

		for typeToImport, _ := range typesToImport {
			res += fmt.Sprintf("%s,\n", typeToImport)
		}

		res += fmt.Sprintf("} from \"%s\"\n", importSource)
	}

	return res
}
