package radicle

type PrimOp func(*State, []Value) (*State, Value, error)

func MapPrimOp(id Ident) PrimOp {
	return nil
	// TODO
}
