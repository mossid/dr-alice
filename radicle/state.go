package radicle

// Core.Bindings
type State struct {
	Env    Env
	PrimFn func(Ident) PrimOp
	//	Mem map[Ref]Value
}

func (s *State) ModifyEnv(f func(Env) Env) *State {
	return &State{
		Env: f(s.Env),
	}
}
