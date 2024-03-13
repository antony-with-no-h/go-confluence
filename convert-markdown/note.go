package convert_markdown

import (
	"bytes"
	"regexp"

	"github.com/gomarkdown/markdown/ast"
)

type Note struct {
	ast.Leaf
	Notes []byte
	Type  []byte
}

func parserNoteHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := parserNote(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func parserNote(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, []byte("!!!")) {
		return nil, nil, 0
	}

	re := regexp.MustCompile(`(?m)!{3}.*$`)
	prefixIndex := re.FindIndex(data)
	prefix := data[prefixIndex[0]:prefixIndex[1]]
	prefixLen := len(prefix)

	re = regexp.MustCompile(`info|note|tip|warning`)
	macroType := re.Find(prefix)

	//prefixLen := len(note)
	endOfBlock := bytes.Index(data[prefixLen:], []byte("\n\n"))
	endOfBlock = endOfBlock + prefixLen

	node := &Note{
		Notes: data[prefixLen:endOfBlock],
		Type:  macroType,
	}

	return node, nil, endOfBlock
}
