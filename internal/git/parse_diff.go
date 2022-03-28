package git

import (
	"bufio"
	"fmt"
	"io"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{reader: bufio.NewReader(reader)}
}

func (parser *Parser) NextDiff() (*AffectedFile, error) {
	c, err := parser.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if c != ':' {
		return nil, fmt.Errorf("expected leading colon in raw diff line")
	}
	d := &AffectedFile{}
	for _, field := range []*string{&d.SrcMode, &d.DstMode, &d.SrcSHA, &d.DstSHA} {
		if *field, err = parser.readStringChop(' '); err != nil {
			return nil, err
		}
	}
	for _, field := range []*string{&d.Status, &d.SrcPath} {
		if *field, err = parser.readStringChop(0); err != nil {
			return nil, err
		}
	}
	if len(d.Status) > 0 && (d.Status[0] == 'C' || d.Status[0] == 'R') {
		if d.DstPath, err = parser.readStringChop(0); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (parser *Parser) readStringChop(delim byte) (string, error) {
	s, err := parser.reader.ReadString(delim)
	if err != nil {
		return "", fmt.Errorf("read raw diff: %v", err)
	}
	return s[:len(s)-1], nil
}
