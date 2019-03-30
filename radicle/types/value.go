package types

type Ident string

func NewIdent(s string) Ident {
	return Ident(s) // TODO
}

func (i Ident) String() string {
	return string(i)
}

type ValueType byte

const (
	VNULL ValueType = iota
	VAtom
	VKeyword
	VString
	VNumber
	VBoolean
	VList
	VVec
	VPrimFn
	VDict
	VRef
	VHandle
	VProcHandle
	VLambda
	VLambdaRec
	VEnv
	VState
)

type Value interface {
	Type() ValueType
}

type Atom Ident

func NewAtom(str string) *Atom { res := Atom(NewIdent(str)); return &res }
func (*Atom) Type() ValueType  { return VAtom }
func (a *Atom) Ident() Ident   { return Ident(*a) }

type Keyword Ident

func NewKeyword(str string) *Keyword { res := Keyword(NewIdent(str)); return &res }
func (k *Keyword) Type() ValueType   { return VKeyword }

type String string

func NewString(str string) *String { res := String(str); return &res }
func (*String) Type() ValueType    { return VString }
func (s *String) String() string   { return string(*s) }

type Bool bool

func NewBool(b bool) *Bool    { res := Bool(b); return &res }
func (*Bool) Type() ValueType { return VBoolean }
func (b *Bool) Bool() bool    { return bool(*b) }

type Num Dec

func NewNum(i int64) *Num    { res := Num(NewDec(i)); return &res }
func (*Num) Type() ValueType { return VNumber }
func (n *Num) Num() Dec      { return Dec(*n) }

type Lambda struct {
	Args   []Ident
	Bodies []Value
	Env    Env
}

func (*Lambda) Type() ValueType { return VLambda }

type LambdaRec struct {
	Self Ident
	*Lambda
}

func (*LambdaRec) Type() ValueType { return VLambdaRec }

type List []Value

func NewList(vs ...Value) *List { res := List(vs); return &res }
func (*List) Type() ValueType   { return VList }
func (l *List) List() []Value   { return []Value(*l) }

type Vector []Value

func NewVector(vs ...Value) *Vector { res := Vector(vs); return &res }
func (*Vector) Type() ValueType     { return VVec }
func (v *Vector) Vector() []Value   { return []Value(*v) }

type PrimFn Ident

func (*PrimFn) Type() ValueType { return VPrimFn }
func (p *PrimFn) Ident() Ident  { return Ident(*p) }

type Dict map[Value]Value

func (*Dict) Type() ValueType                  { return VDict }
func (d *Dict) Get(k Value) (v Value, ok bool) { v, ok = (*d)[k]; return }
func (d *Dict) Set(k, v Value)                 { (*d)[k] = v }

// type Ref *interface{}
