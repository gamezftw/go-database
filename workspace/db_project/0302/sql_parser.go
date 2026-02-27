package db0302

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	buf string
	pos int
}

func NewParser(s string) Parser {
	return Parser{buf: s, pos: 0}
}

func isSpace(ch byte) bool {
	switch ch {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}
func isAlpha(ch byte) bool {
	return 'a' <= (ch|32) && (ch|32) <= 'z'
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isNameStart(ch byte) bool {
	return isAlpha(ch) || ch == '_'
}
func isNameContinue(ch byte) bool {
	return isAlpha(ch) || isDigit(ch) || ch == '_'
}
func isSeparator(ch byte) bool {
	return ch < 128 && !isNameContinue(ch)
}
func isDigitSymbolPrefix(ch byte) bool {
	return ch == '-' || ch == '+'
}
func isQuote(ch byte) bool {
	return ch == '"' || ch == '\''
}
func isBackslash(ch byte) bool {
	return ch == '\\'
}
func isEscapable(ch byte) bool {
	return isQuote(ch) || isBackslash(ch)
}

func (p *Parser) skipSpaces() {
	for p.pos < len(p.buf) && isSpace(p.buf[p.pos]) {
		p.pos += 1
	}
}

func (p *Parser) tryKeyword(kw string) bool {
	p.skipSpaces()
	if !(p.pos+len(kw) <= len(p.buf) && strings.EqualFold(p.buf[p.pos:p.pos+len(kw)], kw)) {
		return false
	}
	if p.pos+len(kw) < len(p.buf) && !isSeparator(p.buf[p.pos+len(kw)]) {
		return false
	}
	p.pos += len(kw)
	return true
}

func (p *Parser) tryName() (string, bool) {
	p.skipSpaces()
	start, cur := p.pos, p.pos
	if !(cur < len(p.buf) && isNameStart(p.buf[cur])) {
		return "", false
	}
	cur++
	for cur < len(p.buf) && isNameContinue(p.buf[cur]) {
		cur++
	}
	p.pos = cur
	return p.buf[start:cur], true
}

func (p *Parser) parseValue(out *Cell) error {
	p.skipSpaces()
	if p.pos >= len(p.buf) {
		return errors.New("expect value")
	}
	ch := p.buf[p.pos]
	if ch == '"' || ch == '\'' {
		return p.parseString(out)
	} else if isDigit(ch) || ch == '-' || ch == '+' {
		return p.parseInt(out)
	} else {
		return errors.New("expect value")
	}
}

func (p *Parser) parseString(out *Cell) error {
	start, cur := p.pos, p.pos
	if !isQuote(p.buf[cur]) {
		return errors.New("invalid string")
	}
	cur += 1 // skip initial quote

	tempStart := start + 1
	str := make([]byte, 0) // TODO: replace this with slice of length of len(p.buf) - p.pos? would this improve performacne since there would be no copying of the slice
	for cur < len(p.buf) {
		curCh := p.buf[cur]
		if isBackslash(curCh) && isEscapable(p.buf[cur+1]) {
			str = append(str, []byte(p.buf[tempStart:cur])...)
			str = append(str, p.buf[cur+1])
			cur += 2
			tempStart = cur
		} else if isQuote(curCh) {
			break
		} else {
			cur += 1
		}
	}
	if !isQuote(p.buf[cur]) {
		return errors.New("missing quotes")
	}

	p.pos = cur + 1
	out.Type = TypeStr
	if len(str) != 0 {
		str = append(str, []byte(p.buf[tempStart:cur])...)
		out.Str = str
	} else {
		out.Str = []byte(p.buf[start+1 : cur])
	}

	return nil
}

func (p *Parser) parseInt(out *Cell) (err error) {
	start, cur := p.pos, p.pos
	ch := p.buf[p.pos]

	fmt.Println(p.buf[start:cur])

	if isDigitSymbolPrefix(ch) {
		cur += 1
	}
	fmt.Println(p.buf[start:cur])

	for cur < len(p.buf) {
		curCh := p.buf[cur]
		if isDigit(curCh) {
			cur += 1
			continue
		} else if isSpace(curCh) {
			break
		} else {
			return errors.New("invalid int value")
		}
	}
	p.pos = cur
	out.Type = TypeI64
	out.I64, err = strconv.ParseInt(p.buf[start:cur], 10, 64)
	return err
}

func (p *Parser) isEnd() bool {
	p.skipSpaces()
	return p.pos >= len(p.buf)
}

// QzBQWVJJOUhU https://trialofcode.org/
