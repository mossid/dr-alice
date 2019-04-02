package radicle

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mossid/dr-alice/parse"
	"github.com/mossid/dr-alice/types"
)

func SetEnv(env types.Env, kvs ...string) Env {
	for i := 0; i < len(kvs); i += 2 {
		env = env.Set(types.NewIdent(kvs[i]), parse.Expr(kvs[i+1]))
	}
	return env
}

func church(body string) string {
	return fmt.Sprintf(`
	(do 
		(def s (fn [n] (fn [f x] (f (n f x)))))
		(def z         (fn [f x] x))
		%s)
	)
	`, body)
}

func parseEval(env types.Env, expr string) types.Value {
	state := EmptyState().SetEnv(env)
	_, v, err := BaseEval(state, parse.Expr(expr))
	if err != nil {
		panic(err)
	}
	return v
}

func TestBaseEval(t *testing.T) {
	env := types.NewListEnv()

	tcs := []struct {
		Env    Env
		Input  Value
		Result Value
	}{
		{
			// 0
			// Basic Atom evaluation
			// atom1 ==> :keyword1
			SetEnv(env, "atom1", ":keyword1"),
			parse.Expr("atom1"),
			parse.Expr(":keyword1"),
		},
		{
			// 1
			// Basic Vector evaluation
			// [atom1, [atom2]] ==> [:keyword1, [keyword2]]
			SetEnv(env, "atom1", ":keyword1", "atom2", ":keyword2"),
			parse.Expr("[atom1 [atom2]]"),
			parse.Expr("[:keyword1 [:keyword2]]"),
		},
		{
			// 2
			// Basic Dict evaluation
			// {:key1 0 :key2 #f}
			SetEnv(env, "atom1", ":keyword1", "atom2", ":keyword2"),
			parse.Expr("{atom1 3 :key [4 atom2]}"),
			parse.Expr("{:keyword1 3 :key [4 :keyword2]}"),
		},
		{
			// 3
			// Def and eval
			env,
			parse.Expr("(do (def x 3) x)"),
			parse.Expr("3"),
		},
		{
			// 4
			// If
			env,
			parse.Expr("(do (def x #f) (if x 3 4))"),
			parse.Expr("4"),
		},
		{
			// 5
			// Lambda
			env,
			parse.Expr("(do (def neg (fn [x] (if x #f #t))) (neg (neg 3)))"),
			parse.Expr("#t"),
		},
		{
			// 6
			// Quote
			env,
			parse.Expr("(do (def l '(3 4 5 6)) l)"),
			parse.Expr("(3 4 5 6)"),
		},
		{
			// 7
			// Church numeric
			env,
			parse.Expr(church("z")),
			parseEval(env, "(fn [f x] x)"),
		},
		{
			// 8
			// PrimFn +
			env,
			parse.Expr("(+ 1 2)"),
			parse.Expr("3"),
		},
		{
			// 9
			// PrimFn used in fn defs
			env,
			parse.Expr("(do (def add (fn [x y] (+ x y))) (add 4 5))"),
			parse.Expr("9"),
		},
		{
			// 10
			// Church numeric evaluated
			env,
			parse.Expr(church(`(do 
				(def add1 (fn [x] (+ x 1)))
				(def x (s (s (s (s (s z))))))
				(x add1 0)
			)`)),
			parse.Expr("5"),
		},
		{
			// 11
			// Ref read
			env,
			parse.Expr(`(do (def r (ref 3)) (read-ref r))`),
			parse.Expr("3"),
		},
		{
			// 12
			// Ref write read
			env,
			parse.Expr(`(do (def r (ref 3)) (write-ref r 4) (read-ref r))`),
			parse.Expr("4"),
		},
	}

	for i, tc := range tcs {
		fmt.Println("Running test", i)
		state := EmptyState().SetEnv(tc.Env)
		_, v, err := BaseEval(state, tc.Input)
		if tc.Result != nil {
			require.NoError(t, err, "Got error(%d): %s", i, err)
			require.True(t, tc.Result.Equal(v), "Not equal(%d): expected %s but %s", i, tc.Result, v)
		} else {
			require.NoError(t, err, "Got success(%d): i, %s")
		}
	}
}
