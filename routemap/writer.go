package routemap

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// ModifiedRoutemap will read the lines of a routemap and insert missing
// content at the releveant locations
func ModifiedRoutemap(file io.Reader, templates, handlers []Missing) ([]byte, error) {
	tplByLine := make(map[int]Missing)
	for i := 0; i < len(templates); i++ {
		tplByLine[templates[i].Position.Line] = templates[i]
	}
	hlrByLine := make(map[int]Missing)
	for i := 0; i < len(handlers); i++ {
		hlrByLine[handlers[i].Position.Line] = handlers[i]
	}
	lineScanner := bufio.NewScanner(file)
	var line int
	var out bytes.Buffer
	for lineScanner.Scan() {
		lineText := lineScanner.Text()
		fmt.Fprintln(&out, lineText)
		line++
		if tmpl, ok := tplByLine[line]; ok {
			fmt.Fprintf(&out, "%s_template = %s\n", getIndentation(tmpl.Position.Col, lineText), strconv.Quote(tmpl.InsertContent))
		}
		if hlr, ok := hlrByLine[line]; ok {
			fmt.Fprintf(&out, "%s_handler = %s\n", getIndentation(hlr.Position.Col, lineText), strconv.Quote(hlr.InsertContent))
		}
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func getIndentation(col int, prevLine string) string {
	if col < len(prevLine) {
		return prevLine[:col-1]
	}
	return prevLine
}
