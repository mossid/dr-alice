package parse

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var yesstrings = []string{`"hello"`, `"world world"`, `""`}
var nostrings = []string{`"pppp`, `"`}
var yesbools = []string{`#t`, `#f`}
var nobools = []string{`#`}
var keywords = []string{`:hello`, `:00123`, `:$@//`}
var nums = []string{`+36`, `-9986`, `324567`}
var atoms = []string{`atom`, `<atom1234@:`}

var values0 = concat(yesstrings, yesbools, keywords, nums, atoms)
var quotes0 = quote(values0)
var lists0 = genlist(values0)
var vecs0 = genvec(values0)
var dicts0 = gendict(values0)

var values1 = concat(values0, quotes0, lists0, vecs0, dicts0)
var quotes1 = quote(values1)
var lists1 = genlist(values1)
var vecs1 = genvec(values1)
var dicts1 = gendict(values1)

var values2 = concat(values1, quotes1, lists1, vecs1, dicts1)
var quotes2 = quote(values2)
var lists2 = genlist(values2)
var vecs2 = genvec(values2)
var dicts2 = gendict(values2)

var values = concat(values2, quotes2, lists2, vecs2, dicts2)

func concat(strss ...[]string) (res []string) {
	for _, strs := range strss {
		res = append(res, strs...)
	}
	return
}

func gen(b, e string, vs []string) (res []string) {
	res = make([]string, 5)
	for n := range res {
		elems := make([]string, 6)
		for i := range elems {
			elems[i] = vs[rand.Intn(len(vs))]
		}
		res[n] = b + strings.Join(elems, " ") + e
	}
	return
}

func genlist(vs []string) []string {
	return gen("(", ")", vs)
}

func genvec(vs []string) []string {
	return gen("[", "]", vs)
}

func gendict(vs []string) []string {
	return gen("{", "}", vs)
}

func quote(vs []string) (res []string) {
	res = make([]string, len(vs))
	for i, v := range vs {
		res[i] = "'" + v
	}
	return
}

func TestValue(t *testing.T) {
	for _, str := range values {
		require.NotNil(t, Value(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range concat(nostrings, nobools) {
		require.Nil(t, Value(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range yesstrings {
		require.NotNil(t, StringLiteral(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range concat(nostrings, yesbools, nobools, keywords, nums, atoms) {
		require.Nil(t, StringLiteral(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestBoolLiteral(t *testing.T) {
	for _, str := range yesbools {
		require.NotNil(t, BoolLiteral(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range concat(yesstrings, nostrings, nobools, keywords, nums, atoms) {
		require.Nil(t, BoolLiteral(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestKeyword(t *testing.T) {
	for _, str := range keywords {
		require.NotNil(t, Keyword(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range concat(yesstrings, nostrings, yesbools, nobools, nums, atoms) {
		require.Nil(t, Keyword(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestNum(t *testing.T) {
	for _, str := range nums {
		require.NotNil(t, NumLiteral(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range concat(yesstrings, nostrings, yesbools, nobools, keywords, atoms) {
		require.Nil(t, NumLiteral(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestAtom(t *testing.T) {
	for _, str := range atoms {
		require.NotNil(t, Atom(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	// Nums can be parsed as atoms
	for _, str := range concat(yesstrings, nostrings, yesbools, nobools, keywords) {
		require.Nil(t, Atom(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestQuote(t *testing.T) {
	for _, str := range quotes1 {
		require.NotNil(t, Quote(&ParserState{[]byte("'" + str)}), "Failed: %s", str)
	}

	for _, str := range concat(values0, lists0, vecs0, lists1, vecs1) {
		require.Nil(t, Quote(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestList(t *testing.T) {
	for _, str := range lists1 {
		require.NotNil(t, List(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range values0 {
		require.Nil(t, List(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}
}

func TestVec(t *testing.T) {
	for _, str := range vecs1 {
		require.NotNil(t, Vec(&ParserState{[]byte(str)}), "Failed: %s", str)
	}

	for _, str := range values0 {
		require.Nil(t, Vec(&ParserState{[]byte(str)}), "Not failed: %s", str)
	}

}

func BenchmarkParsing(b *testing.B) {
	l := 0
	for i := 0; i < b.N; i++ {
		v := values[i%len(values)]
		l += len(v)
		Value(&ParserState{[]byte(v)})
	}

	fmt.Printf("\naverage input length %d for size %d", l/b.N, b.N)
}
