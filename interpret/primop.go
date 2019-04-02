package radicle

import (
	"github.com/mossid/dr-alice/types"
)

type PrimOp func(*State, []Value) (*State, Value, error)

func argn(n int, fn PrimOp) PrimOp {
	return func(s *State, args []Value) (*State, Value, error) {
		if len(args) != n {
			return nil, nil, WrongNumberArgsError("PrimOp", n, len(args))
		}
		return fn(s, args)
	}
}

func MapPrimFn(id Ident) PrimOp {
	switch id.String() {
	case "+":
		return argn(2, func(s *State, args []Value) (*State, Value, error) {
			arg0, ok := args[0].(*Num)
			if !ok {
				return nil, nil, TypeError("+", &Num{}, args[0])
			}
			arg1, ok := args[1].(*Num)
			if !ok {
				return nil, nil, TypeError("+", &Num{}, args[1])
			}
			// TODO: refactor num
			return s, &types.Num{arg0.Add(arg1.Dec)}, nil
		})
	case "ref":
		return argn(1, func(s *State, args []Value) (*State, Value, error) {
			ix := s.Refs.Insert(args[0])
			return s, types.NewRef(ix), nil
		})
	case "read-ref":
		return argn(1, func(s *State, args []Value) (*State, Value, error) {
			arg0, ok := args[0].(*Ref)
			if !ok {
				return nil, nil, TypeError("read-ref", new(types.Ref), args[0])
			}
			res, ok := s.Refs.Get(arg0.Uint())
			if !ok {
				return nil, nil, ImpossibleError("read-ref", "undefined reference")
			}
			return s, res, nil
		})
	case "write-ref":
		return argn(2, func(s *State, args []Value) (*State, Value, error) {
			arg0, ok := args[0].(*Ref)
			if !ok {
				return nil, nil, TypeError("write-ref", new(types.Ref), args[0])
			}
			v := args[1]
			s.Refs.Set(arg0.Uint(), v)
			return s, v, nil
		})
	default:
		return nil
	}
}
