package radicle

import (
	"github.com/mossid/dr-alice/types"
)

const (
	TypeNULL       = types.TypeNULL
	TypeAtom       = types.TypeAtom
	TypeKeyword    = types.TypeKeyword
	TypeString     = types.TypeString
	TypeNumber     = types.TypeNumber
	TypeBoolean    = types.TypeBoolean
	TypeList       = types.TypeList
	TypeVec        = types.TypeVec
	TypePrimFn     = types.TypePrimFn
	TypeDict       = types.TypeDict
	TypeRef        = types.TypeRef
	TypeHandle     = types.TypeHandle
	TypeProcHandle = types.TypeProcHandle
	TypeLambda     = types.TypeLambda
	TypeLambdaRec  = types.TypeLambdaRec
	TypeEnv        = types.TypeEnv
	TypeState      = types.TypeState
)

type (
	ValueType = types.ValueType

	Value  = types.Value
	Ident  = types.Ident
	Env    = types.Env
	Intmap = types.Intmap

	//	Dec = types.Dec

	Atom      = types.Atom
	Bool      = types.Bool
	List      = types.List
	String    = types.String
	Num       = types.Num
	Vector    = types.Vector
	Dict      = types.Dict
	Lambda    = types.Lambda
	LambdaRec = types.LambdaRec
	PrimFn    = types.PrimFn
	Ref       = types.Ref
	State     = types.State
)
