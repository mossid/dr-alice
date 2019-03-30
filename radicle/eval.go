package radicle

func BaseEval(s *State, v Value) (*State, Value, error) {
	switch v := v.(type) {
	case *Atom:
		res, ok := s.Env.Get(v.Ident())
		if !ok {
			return s, v, UnknownIdentifierError(v.Ident())
		}
		return s, res, nil
	case *List:
		l := v.List()
		if len(l) < 1 {
			return nil, nil, WrongNumberArgsError("application", 2, len(l))
		}
		return dollarDollar(s, l[0], l[1:])
	case *Vector:
		xs := v.Vector()
		res := Vector(make([]Value, len(xs)))
		var err error
		for i, x := range xs {
			s, res[i], err = BaseEval(s, x)
			if err != nil {
				return nil, nil, err
			}
		}
		return s, &res, nil
	case *Dict:
		res := Dict(make(map[Value]Value))
		var k0, v0 Value
		var err error
		for k, v := range *v {
			s, k0, err = BaseEval(s, k)
			if err != nil {
				return nil, nil, err
			}
			s, v0, err = BaseEval(s, v)
			if err != nil {
				return nil, nil, err
			}
			res[k0] = v0
		}
		return s, &res, nil
	default:
		return s, v, nil
	}
}

func dollarDollar(s *State, f Value, args []Value) (*State, Value, error) {
	fatom, ok := f.(*Atom)
	if ok {
		f0 := MapSpecialForm(fatom.Ident())
		if f0 != nil {
			return f0(s, args)
		}
	}

	s0, f0, err := BaseEval(s, f)
	if err != nil {
		return nil, nil, err
	}
	args0 := make([]Value, len(args))
	for i, arg := range args {
		s0, args0[i], err = BaseEval(s0, arg)
		if err != nil {
			return nil, nil, err
		}
	}
	return callFn(s0, f0, args0)
}

/*
func Eval(s *State, v Value) (*State, Value, error) {
	e := s.GetEnv(NewIdent("eval"))

}
*/
/*
func EvalExpr(s State, v0 *List) (State, Value, error) {
	v := v0.List()
	fatom, ok := v[0].(*Atom)
	if !ok {
		// TODO: apply instead of error
		return s, v, errors.New("invalid type in list")
	}
	switch fatom.Ident().String() {
	case "cond":
		return
	case "fn":
	case "":
	}
}

func ApplyBase(s State, v Value) (State, Value, error) {

}
*/
