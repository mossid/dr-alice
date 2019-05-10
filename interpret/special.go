package radicle

import "github.com/mossid/dr-alice/types"

type SpecialForm func(*Bindings, []Value) (*Bindings, Value, error)

// TODO: catch & match
func MapSpecialForm(id Ident) SpecialForm {
	switch id {
	case "fn":
		return fn
	case "quote":
		return quote
	case "def":
		return def
	case "def-rec":
		return defrec
	case "do":
		return do
	case "if":
		return iff
	case "cond":
		return cond
	case "module":
		return module
	default:
		return nil
	}
}

func fn(s *Bindings, v []Value) (*Bindings, Value, error) {
	if len(v) < 2 {
		return nil, nil, SpecialFormError("fn", "need an argument vector and a body")
	}
	args, bs := v[0], v[1:]

	vargs, ok := args.(*Vector)
	if !ok {
		return nil, nil, SpecialFormError("fn", "first argument must be a vector of argument atoms")
	}

	iargs := make([]Ident, len(*vargs))

	for i, arg := range vargs.Vector() {
		iarg, ok := arg.(*Atom)
		if !ok {
			return nil, nil, SpecialFormError("fn", "one of the arguments was not an atom")
		}
		iargs[i] = iarg.Ident()
	}

	return s, &Lambda{iargs, bs, s.Env.CloneImmutable()}, nil
}

func quote(s *Bindings, v []Value) (*Bindings, Value, error) {
	if len(v) != 1 {
		return nil, nil, WrongNumberArgsError("quote", 1, len(v))
	}

	return s, v[0], nil
}

func defintern(s *Bindings, v []Value, isrec bool) (*Bindings, Value, error) {
	var fnname string
	if isrec {
		fnname = "def-rec"
	} else {
		fnname = "def"
	}

	if len(v) != 2 {
		return nil, nil, WrongNumberArgsError(fnname, 2, len(v))
	}

	name, ok := v[0].(*Atom)
	if !ok {
		return nil, nil, SpecialFormError(fnname, "expects atom for first arg")
	}
	ident := name.Ident()

	// TODO: implement commenting

	s0, body, err := BaseEval(s, v[1])
	if err != nil {
		return nil, nil, err
	}

	if !isrec {
		s00 := s0.ModifyEnv(func(env Env) Env { return env.Set(ident, body) })
		return s00, nil, nil
	}

	switch body := body.(type) {
	case *Lambda:
		s00 := s0.ModifyEnv(func(env Env) Env {
			return env.Set(ident, &LambdaRec{ident, body})
		})
		return s00, nil, nil
	case *LambdaRec:
		return nil, nil, SpecialFormError(fnname, "cannot be used to alias functions")
	default:
		return nil, nil, SpecialFormError(fnname, "can only be used to define functions")
	}
}

func def(s *Bindings, v []Value) (*Bindings, Value, error) {
	return defintern(s, v, false)
}

func defrec(s *Bindings, v []Value) (*Bindings, Value, error) {
	return defintern(s, v, true)
}

func do(s *Bindings, v []Value) (s0 *Bindings, res Value, err error) {
	s0 = s
	for _, v0 := range v {
		s0, res, err = BaseEval(s0, v0)
	}
	return
}

/*
func catch(s *Bindings, v []Value) (s0 *Bindings, res Value, err error) {

}
*/
func iff(s *Bindings, v []Value) (*Bindings, Value, error) {
	if len(v) != 3 {
		return nil, nil, WrongNumberArgsError("if", 3, len(v))
	}

	s0, cond, err := BaseEval(s, v[0])
	if err != nil {
		return nil, nil, err
	}

	body := v[1]
	// Yes I hate this too
	bcond, ok := cond.(*Bool)
	if ok {
		if bcond.Bool() == false {
			body = v[2]
		}
	}

	return BaseEval(s0, body)
}

func cond(s *Bindings, v []Value) (*Bindings, Value, error) {
	if len(v)%2 != 0 {
		return nil, nil, WrongNumberArgsError("cond", 2, len(v))
	}

	if len(v) == 0 {
		return s, nil, nil
	}

	s0, c, err := BaseEval(s, v[0])
	if err != nil {
		return nil, nil, err
	}

	bcond, ok := c.(*Bool)
	if ok {
		if bcond.Bool() == false {
			return cond(s0, v[2:])
		}
	}

	return BaseEval(s0, v[1])
}

func meta(v Value) (res ModuleMeta, err error) {
	d, ok := v.(*Dict)
	if !ok {
		err = InvalidDeclaration("must be dict", v)
		return
	}
	d0 := *d
	module, ok := d0[types.NewKeyword("module")]
	if !ok {
		err = InvalidDeclaration("missing :module key", v)
		return
	}
	module0, ok := module.(*Atom)
	if !ok {
		err = InvalidDeclaration(":module must be an atom", v)
		return
	}
	doc, ok := d0[types.NewKeyword("doc")]
	if !ok {
		err = InvalidDeclaration("missing :doc key", v)
		return
	}
	doc0, ok := doc.(*String)
	if !ok {
		err = InvalidDeclaration(":doc must be a string", v)
		return
	}
	exports, ok := d0[types.NewKeyword("exports")]
	if !ok {
		err = InvalidDeclaration("missing :exports key", v)
		return
	}
	exports0, ok := exports.(*Vector)
	if !ok {
		err = InvalidDeclaration(":exports must be a vector", v)
		return
	}

	es := make([]string, exports0.Length())
	exports0.Iterate(func(i int, v Value) (abort bool) {
		v0, ok := v.(*Atom)
		if !ok {
			err = InvalidDeclaration(":exports must be a vector of atoms", v)
			return true
		}
		es[i] = v0.Ident()
		return
	})

	return ModuleMeta{module0.Ident(), es, doc0.String()}, nil
}
func module(s *Bindings, v []Value) (*Bindings, Value, error) {
	if len(v) < 1 {
		return nil, nil, WrongNumberArgsError("module", 1, len(v))
	}
	s0, m, err := BaseEval(s, v[0])
	if err != nil {
		return nil, nil, err
	}
	meta, err := meta(m)
	if err != nil {
		return nil, nil, err
	}
	// XXX: make a new scope
	for _, form := range v[1:] {
		s0, _, err = BaseEval(s, form) // XXX: change to Eval
		if err != nil {
			return nil, nil, err
		}
	}

	env := s0.Env

	for _, e := range meta.Exports {
		if _, ok := env.Get(e); !ok {
			return nil, nil, UndefinedExports(e)
		}
	}

	exports := make([]types.Value, len(meta.Exports))
	for i, e := range meta.Exports {
		exports[i] = types.NewAtom(e)
	}

	modu := types.NewDict(
		types.NewKeyword("module"), types.NewAtom(meta.Name),
		types.NewKeyword("env"), env.CloneMutable(),
		types.NewKeyword("exports"), types.NewVector(exports...),
	)

	return s, modu, nil // XXX
}

/*
func match(s *Bindings, v []Value) (*Bindings, Value, error) {
	if len(v) < 1 {
		return nil, nil, PatternMatchError("match", "no value")
	}
	if len(v[1:])%2 != 0 {
		return nil, nil, WrongNumberArgsError("match", 2, len(v))
	}

	s0, v0, err := BaseEval(s, v[0])
	if err != nil {
		return nil, nil, err
	}

	return goMatches(s0, v0, v[1:])
}

func goMatches(s *Bindings, v Value, cases []Value) (*Bindings, Value, error {
	if len(cases) == 0 {
		return nil, nil, PatternMatchError("match", "no match")
	}

	// Inlining match-pat primfn
	// It feels extremely dangerous match-pat is modifiable



}
*/
