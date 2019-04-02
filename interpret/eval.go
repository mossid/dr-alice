package radicle

import (
	"github.com/mossid/dr-alice/types"
)

type Intmap struct {
	m    map[uint64]Value
	next uint64
}

func NewIntmap() *Intmap {
	return &Intmap{
		m:    make(map[uint64]Value),
		next: 0,
	}
}

func (m *Intmap) Insert(v Value) uint64 {
	ix := m.next
	m.m[ix] = v
	m.next = ix + 1
	return ix
}

func (m *Intmap) Get(ix uint64) (res Value, ok bool) {
	res, ok = m.m[ix]
	return
}

func (m *Intmap) Set(ix uint64, v Value) {
	m.m[ix] = v
}

// Core.Bindings
type State struct {
	Env    Env
	PrimFn func(Ident) PrimOp
	Refs   *Intmap
	//	Mem map[Ref]Value
}

func NewState(env Env, fn func(Ident) PrimOp, refs *Intmap) *State {
	return &State{
		Env:    env,
		PrimFn: fn,
		Refs:   refs,
	}
}

func EmptyState() *State {
	return NewState(types.NewListEnv(), MapPrimFn, NewIntmap())
}

func (s *State) ModifyEnv(f func(Env) Env) *State {
	return NewState(f(s.Env), s.PrimFn, s.Refs)
}

func (s *State) SetEnv(env Env) *State {
	return NewState(env, s.PrimFn, s.Refs)
}

func BaseEval(s *State, v Value) (*State, Value, error) {
	switch v := v.(type) {
	case *Atom:
		// BEGIN
		// primop impl is different from the reference
		// I couldn't found how does the primfns are generated from atoms
		primfn := s.PrimFn(v.Ident())
		if primfn != nil {
			return s, types.NewPrimFn(v), nil
		}
		// END

		res, ok := s.Env.Get(v.Ident())
		if ok {
			return s, res, nil
		}
		return s, v, UnknownIdentifierError(v.Ident())
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
		s0 := s.SetEnv(env)
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
