package radicle

type SpecialForm func(*Bindings, []Value) (*Bindings, Value, error)

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

func module(s *Bindings, v []Value) (*Bindings, Value, error) {
	return nil, nil, nil // XXX
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
