package radicle

import (
	"github.com/mossid/dr-alice/types"
)

type PrimOpRun func(*Bindings, []Value) (*Bindings, Value, error)

type PrimOp struct {
	Name string
	Run  PrimOpRun
}

func (fn PrimOp) argn(n int) PrimOp {
	run := func(s *Bindings, args []Value) (*Bindings, Value, error) {
		if len(args) != n {
			return nil, nil, WrongNumberArgsError(fn.Name, n, len(args))
		}
		return fn.Run(s, args)
	}
	return PrimOp{fn.Name, run}
}

func (fn PrimOp) types(tys ...ValueType) PrimOp {
	run := func(s *Bindings, args []Value) (*Bindings, Value, error) {
		for i, ty := range tys {
			if ty != TypeNULL {
				if args[i].Type() != ty {
					return nil, nil, TypeError(fn.Name, ty, args[i].Type())
				}
			}
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
		PrimOp{"base-eval", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg1 := args[1].(*State)
			s0 := s.SetEnv(arg1.Env).SetRefs(arg1.Refs)
			s1, res, err := BaseEval(s0, args[0])
			if err != nil {
				return nil, nil, err
			}
			return s, types.NewList(res, types.NewState(s1.Env, s1.Refs)), nil
		}}.argn(2).types(TypeNULL, TypeState),
		/*
			argn(0, PrimOp{"pure-state", func(s *Bindings, _ []Value) (*Bindings, Value, error) {
				return s, (&Bindings{
					Env: types.NewListEnv().Set("eval", types.NewPrimFn("base-eval")),
				}).ToRadicle(), nil
			}}),
		*/
		PrimOp{"state->env", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return s, args[0].(*State).Env, nil
		}}.argn(1).types(TypeState),
		/*
			PrimOp{"set-binding", func}
			PrimOp{"get-binding"},
			PrimOp{"set-env"},
		*/
		PrimOp{"apply", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return callFn(s, args[0], args[1].(*List).List())
		}}.argn(2).types(TypeNULL, TypeList),
		PrimOp{"list", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return s, types.NewList(args...), nil
		}},
		PrimOp{"dict", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			if len(args)%2 != 0 {
				return nil, nil, WrongNumberArgsError("dict", 2, len(args))
			}
			return s, types.NewDict(args...), nil
		}},
		PrimOp{"throw", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return nil, nil, ThrownError(args[0].(*Atom).Ident(), args[1])
		}}.argn(2).types(TypeAtom),
		PrimOp{"eq?", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			return s, types.NewBool(args[0].Equal(args[1])), nil
		}}.argn(2),
		PrimOp{"add-right", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			res := append(*args[0].(*Vector), args[1])
			return s, &res, nil
		}}.argn(2).types(TypeVec),
		// Lists
		PrimOp{"cons", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch tail := args[1].(type) {
			case *List:
				return s, &types.List{Head: args[0], Tail: tail}, nil
			case *Vector:
				res := types.Vector(append([]Value{args[0]}, tail.Vector()...))
				return s, &res, nil
			default:
				return nil, nil, TypeError("cons", TypeList, args[1].Type())
			}
		}}.argn(2),
		PrimOp{"first", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch list := args[0].(type) {
			case *List:
				if list == nil {
					return nil, nil, OtherError("first", "Empty list")
				}
				return s, list.Head, nil
			case *Vector:
				if len(*list) == 0 {
					return nil, nil, OtherError("first", "Empty vector")
				}
				return s, (*list)[0], nil
			default:
				return nil, nil, TypeError("first", TypeList, args[0].Type())
			}
		}}.argn(1),
		PrimOp{"rest", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch list := args[0].(type) {
			case *List:
				if list == nil {
					return nil, nil, OtherError("rest", "Empty List")
				}
				return s, list.Tail, nil
			case *Vector:
				if len(*list) == 0 {
					return nil, nil, OtherError("rest", "Empty vector")
				}
				res := (*list)[1:]
				return s, &res, nil
			default:
				return nil, nil, TypeError("rest", TypeList, args[0].Type())
			}
		}}.argn(1),
		// Sequences
		PrimOp{"length", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg0, ok := args[0].(types.Sequence)
			if !ok {
				return nil, nil, TypeError("length", TypeList, args[0].Type())
			}
			return s, types.NewNum(int64(arg0.Length())), nil
		}}.argn(1),
		PrimOp{"drop", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg1, ok := args[1].(types.Sequence)
			if !ok {
				return nil, nil, TypeError("drop", TypeList, args[1].Type())
			}
			return s, arg1.Slice(int(args[0].(*Num).Num()), -1), nil
		}}.argn(2).types(TypeNumber),
		PrimOp{"take", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			arg1, ok := args[1].(types.Sequence)
			if !ok {
				return nil, nil, TypeError("take", TypeList, args[1].Type())
			}
			return s, arg1.Slice(0, int(args[0].(*Num).Num())), nil
		}}.argn(2).types(TypeNumber),
		PrimOp{"nth", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			switch list := args[1].(type) {

			case *List, *Vector:
				return s, list.(types.Sequence).Index(int(args[0].(*Num).Num())), nil
			default:
				return nil, nil, TypeError("nth", TypeList, args[1].Type())
			}
		}}.argn(2).types(TypeNumber),
		// PrimOp{"sort-by"}
		// PrimOp{"zip"}
		// PrimOp{"vec-to-list"}
		// PrimOp{"list-to-vec"}
		// Dicts
		PrimOp{"lookup", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			res, ok := (*args[1].(*Dict))[args[0]]
			if !ok {
				return nil, nil, OtherError("lookup", "key did not exist: " /*, args[0]*/)
			}
			return s, res, nil
		}}.argn(2).types(TypeNULL, TypeDict),
		/*
			PrimOp{"map-values", func(s *Bindings, args []Value) (*Bindings, Value, error) {
				res := Dict(make(map[Value]Value))
				for k, v := range *args[1].(*Dict) {
					res[k] =
				}
			}}.argn(2).types(TypeNULL, TypeDict),
		*/
		// Ref
		PrimOp{"ref", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			ix := s.Refs.Insert(args[0])
			return s, types.NewRef(ix), nil
		}}.argn(1),
		PrimOp{"read-ref", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			res, ok := s.Refs.Get(args[0].(*Ref).Uint())
			if !ok {
				return nil, nil, ImpossibleError("read-ref", "undefined reference")
			}
			return s, res, nil
		}}.argn(1).types(TypeRef),
		PrimOp{"write-ref", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			v := args[1]
			s.Refs.Set(args[0].(*Ref).Uint(), v)
			return s, v, nil
		}}.argn(2).types(TypeRef),
		PrimOp{"+", func(s *Bindings, args []Value) (*Bindings, Value, error) {
			// TODO: refactor num
			return s, types.NewNum(int64(*args[0].(*Num) + *args[1].(*Num))), nil
		}}.argn(2).types(TypeNumber, TypeNumber),
	}
}
