package parse

import (
	"bytes"
	"regexp"
)

type ParserState struct {
	Stream []byte // io.Reader?
}

func (st *ParserState) Lookahead(n int) string {
	return string(st.Stream[:n])
}

func (st *ParserState) CheckString(s string) bool {
	bz := []byte(s)
	if len(st.Stream) < len(bz) {
		return false
	}
	if !bytes.Equal(st.Stream[:len(bz)], bz) {
		return false
	}
	return true
}

func (st *ParserState) Check(regex *regexp.Regexp) (string, bool) {
	loc := regex.FindIndex(st.Stream)
	if loc == nil {
		return "", false
	}
	if loc[0] != 0 {
		return "", false
	}
	return string(st.Stream[:loc[1]]), true
}

// CONTRACT: Check() before
func (st *ParserState) Consume(n int) {
	st.Stream = st.Stream[n:]
}

func (st *ParserState) CheckConsumeString(s string) (ok bool) {
	ok = st.CheckString(s)
	if !ok {
		return
	}
	st.Consume(len(s))
	return
}

func (st *ParserState) CheckConsume(regex *regexp.Regexp) (res string, ok bool) {
	res, ok = st.Check(regex)
	if !ok {
		return
	}
	st.Consume(len(res))
	return
}

func (st *ParserState) CheckConsumeStringEmpty(s string) interface{} {
	if !st.CheckConsumeString(s) {
		return nil
	}
	return Empty{}
}

func (st *ParserState) CheckConsumeEmpty(regex *regexp.Regexp) interface{} {
	_, ok := st.CheckConsume(regex)
	if !ok {
		return nil
	}
	return Empty{}
}

// nil = fail
// Empty{} = no return, success
// any other type = return, success
type Parser func(*ParserState) interface{}

type Empty struct{}

func Try(p Parser) Parser {
	return func(st *ParserState) (res interface{}) {
		orig := *st
		res = p(st)
		if res == nil {
			*st = orig
		}
		return
	}
}

func SkipMany(p Parser) Parser {
	return func(st *ParserState) interface{} {
		for p(st) != nil {
		}
		return Empty{}
	}
}

func Choice(ps ...Parser) Parser {
	return func(st *ParserState) interface{} {
		for _, p := range ps {
			if res := p(st); res != nil {
				return res
			}
		}
		return nil
	}
}

func Space(spc Parser, lcp Parser, bc Parser) Parser {
	return SkipMany(Choice(spc, lcp, bc))
}

func Lexeme(spc Parser, p Parser) Parser {
	return func(st *ParserState) (res interface{}) {
		res = p(st)
		spc(st)
		return
	}
}

func SkipLineComment(pref string) Parser {
	regex := regexp.MustCompile("^" + pref + ".*$")
	return func(st *ParserState) interface{} {
		return st.CheckConsumeEmpty(regex)
	}
}

func SkipBlockComment(start, end string) Parser {
	regex := regexp.MustCompile("^" + start + ".*" + end)
	return func(st *ParserState) interface{} {
		return st.CheckConsumeEmpty(regex)
	}
}

var space1 = regexp.MustCompile(`^[ \n\r\t]`)

func Space1(st *ParserState) interface{} {
	return st.CheckConsumeEmpty(space1)
}
