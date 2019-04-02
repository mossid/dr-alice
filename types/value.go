package types

import (
	"strconv"
	"strings"
)

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
	String() string
	Equal(Value) bool
}

type Atom Ident

func NewAtom(str string) *Atom { res := Atom(NewIdent(str)); return &res }
func (*Atom) Type() ValueType  { return VAtom }
func (a *Atom) Ident() Ident   { return Ident(*a) }
func (a *Atom) String() string { return a.Ident().String() }
func (a *Atom) Equal(v Value) bool {
	a0, ok := v.(*Atom)
	if !ok {
		return false
	}
	return a0.String() == a.String()
}

type Keyword Ident

func NewKeyword(str string) *Keyword { res := Keyword(NewIdent(str)); return &res }
func (k *Keyword) Type() ValueType   { return VKeyword }
func (k *Keyword) Ident() Ident      { return Ident(*k) }
func (k *Keyword) String() string    { return ":" + k.Ident().String() }
func (k *Keyword) Equal(v Value) bool {
	k0, ok := v.(*Keyword)
	if !ok {
		return false
	}
	return k0.String() == k.String()
}

type String string

func NewString(str string) *String { res := String(str); return &res }
func (*String) Type() ValueType    { return VString }
func (s *String) String() string   { return string(*s) } // TODO: fix
func (s *String) Equal(v Value) bool {
	s0, ok := v.(*String)
	if !ok {
		return false
	}
	return *s == *s0
}

type Bool bool

func NewBool(b bool) *Bool    { res := Bool(b); return &res }
func (*Bool) Type() ValueType { return VBoolean }
func (b *Bool) Bool() bool    { return bool(*b) }
func (b *Bool) String() string {
	if *b {
		return "#t"
	}
	return "#f"
}
func (b *Bool) Equal(v Value) bool {
	b0, ok := v.(*Bool)
	if !ok {
		return false
	}
	return *b0 == *b
}

type Num struct {
	Dec
}

func NewNum(i int64) *Num     { return &Num{NewDec(i)} }
func (*Num) Type() ValueType  { return VNumber }
func (n *Num) String() string { return n.String() }
func (n *Num) Equal(v Value) bool {
	n0, ok := v.(*Num)
	if !ok {
		return false
	}
	return n.Dec.Equal(n0.Dec)
}

type Lambda struct {
	Args   []Ident
	Bodies []Value
	Env    Env
}

func (*Lambda) Type() ValueType { return VLambda }
func (l *Lambda) String() string {
	args := make([]string, len(l.Args))
	for i, arg := range l.Args {
		args[i] = arg.String()
	}
	bodies := make([]string, len(l.Bodies))
	for i, body := range l.Bodies {
		bodies[i] = body.String()
	}
	return "(fn [" + strings.Join(args, " ") + "] " + strings.Join(bodies, " ") + ")"
}
func (l *Lambda) Equal(v Value) bool {
	l0, ok := v.(*Lambda)
	if !ok {
		return false
	}
	if len(l0.Args) != len(l.Args) || len(l0.Bodies) != len(l.Bodies) { // TODO: add env
		return false
	}
	for i, arg := range l.Args {
		if arg != l0.Args[i] {
			return false
		}
	}
	for i, body := range l.Bodies {
		if !body.Equal(l0.Bodies[i]) {
			return false
		}
	}
	return true
}

type LambdaRec struct {
	Self Ident
	*Lambda
}

func (l *LambdaRec) Equal(v Value) bool {
	l0, ok := v.(*LambdaRec)
	if !ok {
		return false
	}
	if l0.Self != l.Self {
		return false
	}
	return l.Lambda.Equal(l0.Lambda)
}

func (*LambdaRec) Type() ValueType { return VLambdaRec }

type List []Value

func NewList(vs ...Value) *List { res := List(vs); return &res }
func (*List) Type() ValueType   { return VList }
func (l *List) List() []Value   { return []Value(*l) }
func (l *List) String() string {
	ll := l.List()
	vs := make([]string, len(ll))
	for i, v := range ll {
		vs[i] = v.String()
	}
	return "(" + strings.Join(vs, " ") + ")"
}
func (l *List) Equal(v Value) bool {
	l0, ok := v.(*List)
	if !ok {
		return false
	}
	ll, l0l := l.List(), l0.List()
	if len(ll) != len(l0l) {
		return false
	}
	for i, v := range ll {
		if !l0l[i].Equal(v) {
			return false
		}
	}
	return true
}

type Vector []Value

func NewVector(vs ...Value) *Vector { res := Vector(vs); return &res }
func (*Vector) Type() ValueType     { return VVec }
func (v *Vector) Vector() []Value   { return []Value(*v) }
func (l *Vector) String() string {
	vv := l.Vector()
	vs := make([]string, len(vv))
	for i, v := range vv {
		vs[i] = v.String()
	}
	return "[" + strings.Join(vs, " ") + "]"
}
func (l *Vector) Equal(v Value) bool {
	l0, ok := v.(*Vector)
	if !ok {
		return false
	}
	ll, l0l := l.Vector(), l0.Vector()
	if len(ll) != len(l0l) {
		return false
	}
	for i, v := range ll {
		if !l0l[i].Equal(v) {
			return false
		}
	}
	return true
}

// TODO: mv string hash
type Dict map[Value]Value

func NewDict(kvs ...Value) *Dict {
	if len(kvs)%2 != 0 {
		panic("odd number of arguments in NewDict()")
	}
	res := Dict(make(map[Value]Value))
	for i := 0; i < len(kvs); i += 2 {
		res[kvs[i]] = kvs[i+1]
	}
	return &res
}
func (*Dict) Type() ValueType                  { return VDict }
func (d *Dict) Get(k Value) (v Value, ok bool) { v, ok = (*d)[k]; return }
func (d *Dict) Set(k Value, v Value)           { (*d)[k] = v }
func (d *Dict) String() string {
	var elems []string
	for k, v := range *d {
		elems = append(elems, k.String(), v.String())
	}
	return "{" + strings.Join(elems, " ") + "}"
}
func (d *Dict) Equal(v Value) bool {
	d0, ok := v.(*Dict)
	if !ok {
		return false
	}
	if len(*d) != len(*d0) {
		return false
	}
	strmap := make(map[string]Value)
	for k, v := range *d {
		strmap[k.String()] = v
	}
	for k, v := range *d0 {
		if !strmap[k.String()].Equal(v) {
			return false
		}
	}
	return true
}

type PrimFn Ident

func NewPrimFn(a *Atom) *PrimFn  { res := PrimFn(a.Ident()); return &res }
func (*PrimFn) Type() ValueType  { return VPrimFn }
func (a *PrimFn) Ident() Ident   { return Ident(*a) }
func (a *PrimFn) String() string { return a.Ident().String() }
func (a *PrimFn) Equal(v Value) bool {
	a0, ok := v.(*PrimFn)
	if !ok {
		return false
	}
	return a0.String() == a.String()
}

type Ref uint64

func NewRef(u uint64) *Ref    { res := Ref(u); return &res }
func (*Ref) Type() ValueType  { return VRef }
func (r *Ref) Uint() uint64   { return uint64(*r) }
func (r *Ref) String() string { return "#" + strconv.FormatUint(r.Uint(), 10) }
func (r *Ref) Equal(v Value) bool {
	r0, ok := v.(*Ref)
	if !ok {
		return false
	}
	return r0.Uint() == r.Uint()
}