package parse

import (
	"regexp"
	"strconv"

	"github.com/mossid/dr-alice/radicle/types"
)

func Value(st *ParserState) interface{} {
	v, ok := Choice(
		StringLiteral,
		BoolLiteral,
		Keyword,
		NumLiteral,
		Atom,
		Quote,
		List,
		Vec,
	//	Dict,
	)(st).(types.Value)
	if !ok {
		return nil
	}

	spaceConsume(st)

	return v
}

var skipLineComment = SkipLineComment(`;;`)
var skipBlockComment = SkipBlockComment(`#\|`, `\|#`)

func spaceConsume(st *ParserState) interface{} {
	return Space(Space1, skipLineComment, skipBlockComment)(st)
}

var stringLiteralMatch = regexp.MustCompile(`^"[a-zA-Z0-9 ]*"`)

func StringLiteral(st *ParserState) interface{} {
	// XXX: escape
	str, ok := st.CheckConsume(stringLiteralMatch)
	if !ok {
		return nil
	}

	spaceConsume(st)

	return types.NewString(str[1 : len(str)-1])
}

var boolLiteralMatch = regexp.MustCompile(`^#(t|f)`)

func BoolLiteral(st *ParserState) (res interface{}) {
	if st.CheckConsumeStringEmpty("#t") != nil {
		spaceConsume(st)
		return types.NewBool(true)
	}
	if st.CheckConsumeStringEmpty("#f") != nil {
		spaceConsume(st)
		return types.NewBool(false)
	}
	return nil
}

func Keyword(st *ParserState) interface{} {
	if st.CheckConsumeStringEmpty(":") == nil {
		return nil
	}

	kw, ok := identRest(st).(string)
	if !ok {
		return nil
	}

	spaceConsume(st)

	return types.NewKeyword(kw)
}

var identFirstMatch = regexp.MustCompile(`^[A-Za-z!$%&*+-./<=>?@^_~]`)

func identFirst(st *ParserState) interface{} {
	str, ok := st.CheckConsume(identFirstMatch)
	if !ok {
		return nil
	}
	return str
}

var identRestMatch = regexp.MustCompile(`^[A-Za-z0-9!$%&*+-./:<=>?@^_~]*`)

func identRest(st *ParserState) interface{} {
	str, ok := st.CheckConsume(identRestMatch)
	if !ok {
		return nil
	}
	return str
}

var numLiteralMatch = regexp.MustCompile(`^[+-]?[0-9]+`)

func NumLiteral(st *ParserState) interface{} {

	// TODO: points and rationals
	str, ok := st.CheckConsume(numLiteralMatch)
	if !ok {

		return nil
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {

		return nil
	}

	spaceConsume(st)

	return types.NewNum(i)
}

func Atom(st *ParserState) interface{} {
	l, ok := identFirst(st).(string)
	if !ok {
		return nil
	}

	r, ok := identRest(st).(string)
	if !ok {
		return nil
	}

	spaceConsume(st)

	return types.NewAtom(l + r)
}

func Quote(st *ParserState) interface{} {
	if st.CheckConsumeStringEmpty("'") == nil {
		return nil
	}

	val, ok := Value(st).(types.Value)
	if !ok {
		return nil
	}

	return types.NewList(types.NewAtom("quote"), val)
}

func List(st *ParserState) interface{} {

	if st.CheckConsumeStringEmpty("(") == nil {

		return nil
	}

	spaceConsume(st)

	var vs []types.Value
	for {
		if st.CheckConsumeStringEmpty(")") != nil {

			return types.NewList(vs...)
		}

		v, ok := Value(st).(types.Value)
		if !ok {
			return nil
		}

		vs = append(vs, v)
		spaceConsume(st)
	}
}

func Vec(st *ParserState) interface{} {
	if st.CheckConsumeStringEmpty("[") == nil {
		return nil
	}

	spaceConsume(st)

	var vs []types.Value
	for {
		v, ok := Value(st).(types.Value)
		if !ok {
			break
		}
		vs = append(vs, v)
		spaceConsume(st)
	}

	if st.CheckConsumeStringEmpty("]") == nil {
		return nil
	}

	return types.NewVector(vs...)
}

/*
// TODO
func Dict(st *ParserState) interface{} {

}
*/
