package radicle

func callFn(s *State, v Value, args []Value) (*State, Value, error) {
	switch v := v.(type) {
	case *Lambda:
		if len(v.Args) != len(args) {
			return nil, nil, WrongNumberArgsError("lambda", len(v.Args), len(args))
		}
		env := v.Env
		for i, name := range v.Args {
			env = env.Set(name, args[i])
		}
		s0 := &State{Env: env}
		var res Value
		var err error
		for _, expr := range v.Bodies {
			s0, res, err = BaseEval(s0, expr)
			if err != nil {
				return nil, nil, err
			}
		}
		return s, res, err
	case *LambdaRec:
		l := &Lambda{v.Args, v.Bodies, v.Env.Set(v.Self, v)}
		_, res, err := callFn(nil, l, args)
		return s, res, err
	case *PrimFn:
		fn := s.PrimFn(v.Ident())
		return fn(s, args)
	default:
		return nil, nil, NonFunctionCalledError(v)
	}
}
