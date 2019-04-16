package radicle

import (
	"github.com/mossid/dr-alice/types"
)

type PrimOpRun func(*Bindings, []Value) (*Bindings, Value, error)

type PrimOp struct {
	Name string
	Run  PrimOpRun
}

func argn(n int, fn PrimOp) PrimOp {
	run := func(s *Bindings, args []Value) (*Bindings, Value, error) {
		if len(args) != n {
			return nil, nil, WrongNumberArgsError(fn.Name, n, len(args))
		}
		return fn.Run(s, args)
	}
	return PrimOp{fn.Name, run}
}

func MapPrimOpRuns(ops []PrimOp) func(Ident) PrimOpRun {
	m := make(map[Ident]PrimOpRun)
	for _, op := range ops {
		m[types.NewIdent(op.Name)] = op.Run
	}
	return func(id Ident) PrimOpRun {
		return m[id]
	}
}

func PurePrimFns() []PrimOp {
	return []PrimOp{
		argn(2, PrimOp{"base-eval", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg1, ok := args[1].(*State)
			if !ok {
				return nil, nil, TypeError("base-eval", new(State), args[1])
			}
			s0 := s.SetEnv(arg1.Env).SetRefs(arg1.Refs)
			s1, res, err := BaseEval(s0, args[0])
			if err != nil {
				return nil, nil, err
			}
			return s, types.NewList(res, types.NewState(s1.Env, s1.Refs)), nil
		}}),
		/*
			argn(0, PrimOp{"pure-state", func(s *Bindings, _ []Value) (*Bindings, Value, error) {
				return s, (&Bindings{
					Env: types.NewListEnv().Set("eval", types.NewPrimFn("base-eval")),
				}).ToRadicle(), nil
			}}),
		*/
		argn(2, PrimOp{"apply", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg1, ok := args[1].(*types.List)
			if !ok {
				return nil, nil, TypeError("apply", new(types.List), args[1])
			}
			return callFn(s, args[0], arg1.List())
		}}),
		PrimOp{"list", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return s, types.NewList(args...), nil
		}},
		/*
			PrimOp{"dict", func(s *Bindings, args []Value) (*Bindings, Value, error) {
				if len(args)
			}}
		*/
		argn(2, PrimOp{"eq?", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return s, types.NewBool(args[0].Equal(args[1])), nil
		}}),
		// Lists
		argn(2, PrimOp{"cons", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch tail := args[1].(type) {
			case *List:
				return s, &types.List{Head: args[0], Tail: tail}, nil
			/*
				case *Vector:
					return s, &types.Vector(append(args[0]))
			*/
			default:
				return nil, nil, TypeError("cons", new(types.List), args[1])
			}
		}}),
		argn(1, PrimOp{"first", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch list := args[0].(type) {
			case *List:
				if list == nil {
					return nil, nil, OtherError("first", "Empty list")
				}
				return s, list.Head, nil
			// case *Vector:
			default:
				return nil, nil, TypeError("first", new(types.List), args[0])
			}
		}}),
		argn(1, PrimOp{"rest", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch list := args[0].(type) {
			case *List:
				if list == nil {
					return nil, nil, OtherError("rest", "Empty List")
				}
				return s, list.Tail, nil
				// case *Vector:
			default:
				return nil, nil, TypeError("rest", new(types.List), args[0])
			}
		}}),
		// Sequences
		argn(1, PrimOp{"length", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			cnt := 0
			switch list := args[0].(type) {
			case *List:
				for ; list != nil; list = list.Tail {
					cnt++
				}
			case *Vector:
				cnt = len(list.Vector())
			case *String:
				cnt = len(list.String())
			default:
				return nil, nil, TypeError("length", new(types.List), args[0])
			}
			return s, types.NewNum(int64(cnt)), nil
		}}),
		// Ref
		argn(1, PrimOp{"ref", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			ix := s.Refs.Insert(args[0])
			return s, types.NewRef(ix), nil
		}}),
		argn(1, PrimOp{"read-ref", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg0, ok := args[0].(*Ref)
			if !ok {
				return nil, nil, TypeError("read-ref", new(types.Ref), args[0])
			}
			res, ok := s.Refs.Get(arg0.Uint())
			if !ok {
				return nil, nil, ImpossibleError("read-ref", "undefined reference")
			}
			return s, res, nil
		}}),
		argn(2, PrimOp{"write-ref", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg0, ok := args[0].(*Ref)
			if !ok {
				return nil, nil, TypeError("write-ref", new(types.Ref), args[0])
			}
			v := args[1]
			s.Refs.Set(arg0.Uint(), v)
			return s, v, nil
		}}),
		argn(2, PrimOp{"+", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg0, ok := args[0].(*Num)
			if !ok {
				return nil, nil, TypeError("+", new(Num), args[0])
			}
			arg1, ok := args[1].(*Num)
			if !ok {
				return nil, nil, TypeError("+", new(Num), args[1])
			}
			// TODO: refactor num
			return s, types.NewNum(int64(*arg0 + *arg1)), nil
		}}),
	}
}
