package types

import (
	"strconv"
	"strings"

	"github.com/mossid/dr-alice/types/proto"
)

type Ident = string

func NewIdent(s string) string {
	return s
}

type ValueType byte

const (
	TypeNULL ValueType = iota
	TypeAtom
	TypeKeyword
	TypeString
	TypeNumber
	TypeBoolean
	TypeList
	TypeVec
	TypePrimFn
	TypeDict
	TypeRef
	TypeHandle
	TypeProcHandle
	TypeLambda
	TypeLambdaRec
	TypeEnv
	TypeState
)

type Value interface {
	Type() ValueType
	String() string
	Equal(Value) bool
	Proto() *proto.Value
}

type Atom Ident

func NewAtom(str string) *Atom { res := Atom(NewIdent(str)); return &res }
func (*Atom) Type() ValueType  { return TypeAtom }
func (a *Atom) Ident() Ident   { return Ident(*a) }
func (a *Atom) String() string { return a.Ident() }
func (a *Atom) Equal(v Value) bool {
	a0, ok := v.(*Atom)
	if !ok {
		return false
	}
	return a0.String() == a.String()
}
func (a *Atom) Proto() *proto.Value {
	return &proto.Value{&proto.Value_Atom{&proto.Atom{a.String()}}}
}

type Keyword Ident

func NewKeyword(str string) *Keyword { res := Keyword(NewIdent(str)); return &res }
func (k *Keyword) Type() ValueType   { return TypeKeyword }
func (k *Keyword) Ident() Ident      { return Ident(*k) }
func (k *Keyword) String() string    { return ":" + k.Ident() }
func (k *Keyword) Equal(v Value) bool {
	k0, ok := v.(*Keyword)
	if !ok {
		return false
	}
	return k0.String() == k.String()
}
func (k *Keyword) Proto() *proto.Value {
	return &proto.Value{&proto.Value_Keyword{&proto.Keyword{k.String()}}}
}

type String string

func NewString(str string) *String { res := String(str); return &res }
func (*String) Type() ValueType    { return TypeString }
func (s *String) String() string   { return string(*s) } // TODO: fix
func (s *String) Equal(v Value) bool {
	s0, ok := v.(*String)
	if !ok {
		return false
	}
	return *s == *s0
}
func (s *String) Proto() *proto.Value {
	return &proto.Value{&proto.Value_String_{&proto.String{s.String()}}}
}

func (s *String) Slice(begin, end int) Sequence {
	if end == -1 {
		end = len(*s)
	}
	res := (*s)[begin:end]
	return &res
}

func (s *String) Index(ix int) Value {
	panic("unused function")
}

func (s *String) Length() int {
	return len(*s)
}

func (s *String) Iterate(f func(int, Value) bool) {
	panic("unused function")
}

type Bool bool

func NewBool(b bool) *Bool    { res := Bool(b); return &res }
func (*Bool) Type() ValueType { return TypeBoolean }
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
func (b *Bool) Proto() *proto.Value {
	return &proto.Value{&proto.Value_Boolean{&proto.Boolean{b.Bool()}}}
}
func (b *Bool) Unproto(pv *proto.Value) {
	pa := pv.GetBoolean()
	if pa == nil {
		return
	}
	*b = Bool(pa.Boolean)
}

type Num int64

func NewNum(i int64) *Num     { res := Num(i); return &res }
func (*Num) Type() ValueType  { return TypeNumber }
func (n *Num) String() string { return strconv.FormatInt(n.Num(), 10) }
func (n *Num) Num() int64     { return int64(*n) }
func (n *Num) Equal(v Value) bool {
	n0, ok := v.(*Num)
	if !ok {
		return false
	}
	return n0 == n
}
func (n *Num) Proto() *proto.Value {
	return &proto.Value{&proto.Value_Num{&proto.Num{n.Num()}}}
}
func (n *Num) Unproto(pv *proto.Value) {
	pa := pv.GetNum()
	if pa == nil {
		return
	}
	*n = Num(pa.Num)
}

type Lambda struct {
	Args   []Ident
	Bodies []Value
	Env    Env
}

func NewLambda(args []Ident, bodies []Value, env Env) *Lambda {
	res := Lambda{args, bodies, env}
	return &res
}
func (*Lambda) Type() ValueType { return TypeLambda }
func (l *Lambda) String() string {
	bodies := make([]string, len(l.Bodies))
	for i, body := range l.Bodies {
		bodies[i] = body.String()
	}
	return "(fn [" + strings.Join(l.Args, " ") + "] " + strings.Join(bodies, " ") + ")"
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
func (l *Lambda) Proto() *proto.Value {
	bodies := make([]*proto.Value, len(l.Bodies))
	for i, body := range l.Bodies {
		bodies[i] = body.Proto()
	}
	return &proto.Value{&proto.Value_Lambda{&proto.Lambda{
		l.Args,
		bodies,
		l.Env.Proto().GetEnv(),
	}}}
}
func (l *Lambda) Unproto(pv *proto.Value) {
	/*
		pa := pv.GetLambda()
		if pa == nil {
			return
		}
		bodies := make([]Value, len(pa.Bodies))
		for i, body := range pa.Bodies {
			bodies[i]
		}
		*l = Lambda{
			Args:   pa.Args,
			Bodies: pa.Bodies,
			Env:    pa.Env,
		}
	*/
}

type LambdaRec struct {
	Self Ident
	*Lambda
}

func NewLambdaRec(self Ident, lambda *Lambda) *LambdaRec {
	res := LambdaRec{self, lambda}
	return &res
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

func (*LambdaRec) Type() ValueType { return TypeLambdaRec }

type Sequence interface {
	Value
	// Negative end means empty index: a[3:] -> a.Slice(3, -1)
	Slice(int, int) Sequence
	Index(int) Value
	Length() int
	Iterate(func(int, Value) bool)
}

type List struct {
	Head Value
	Tail *List
}

func NewList(vs ...Value) *List {
	var top *List
	for i := len(vs) - 1; i >= 0; i-- {
		top = &List{
			Head: vs[i],
			Tail: top,
		}
	}
	return top
}
func (*List) Type() ValueType { return TypeList }
func (l *List) List() (res []Value) {
	for ; l != nil; l = l.Tail {
		res = append(res, l.Head)
	}
	return
}
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
func (l *List) Proto() *proto.Value {
	var vs []*proto.Value
	for ; l != nil; l = l.Tail {
		vs = append(vs, l.Head.Proto())
	}
	return &proto.Value{&proto.Value_List{&proto.List{vs}}}
}

func (l *List) Slice(begin, end int) Sequence {
	if l == nil {
		return nil
	}
	if end == 0 {
		return nil
	}

	// Slice all
	if begin == 0 && end == -1 {
		return l
	}

	// Drop
	if end == -1 {
		return l.Tail.Slice(begin-1, end)
	}

	// Take
	if begin == 0 {
		return &List{
			Head: l.Head,
			Tail: l.Tail.Slice(0, end-1).(*List),
		}
	}

	// Normal slicing -> reduce to take
	return l.Tail.Slice(begin-1, end-1)
}

func (l *List) Index(ix int) Value {
	if ix == 0 {
		return l.Head
	}
	return l.Tail.Index(ix - 1)
}

func (l *List) Length() int {
	if l == nil {
		return 0
	}
	return l.Tail.Length() + 1
}

func (l *List) Iterate(f func(ix int, v Value) bool) {
	for i := 0; l != nil; i++ {
		if f(i, l.Head) {
			return
		}
		l = l.Tail
	}
}

// TODO: store in reversed order so we can prepend more efficient
type Vector []Value

func NewVector(vs ...Value) *Vector { res := Vector(vs); return &res }
func (*Vector) Type() ValueType     { return TypeVec }
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
func (v *Vector) Proto() *proto.Value {
	vs := make([]*proto.Value, len(*v))
	for i, v := range *v {
		vs[i] = v.Proto()
	}
	return &proto.Value{&proto.Value_Vector{&proto.Vector{vs}}}
}

func (v *Vector) Slice(begin, end int) Sequence {
	if end == -1 {
		end = len(*v)
	}
	res := (*v)[begin:end]
	return &res
}

func (v *Vector) Index(ix int) Value {
	return (*v)[ix]
}

func (v *Vector) Length() int {
	return len(*v)
}

func (v *Vector) Iterate(f func(int, Value) bool) {
	for i, v := range *v {
		if f(i, v) {
			return
		}
	}
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
func (*Dict) Type() ValueType                  { return TypeDict }
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

func (d *Dict) Proto() *proto.Value {
	// XXX: sort kvpairs, translate to proto.KTypePair
	panic("not implemented")
}

type PrimFn Ident

func NewPrimFn(a *Atom) *PrimFn  { res := PrimFn(a.Ident()); return &res }
func (*PrimFn) Type() ValueType  { return TypePrimFn }
func (a *PrimFn) Ident() Ident   { return Ident(*a) }
func (a *PrimFn) String() string { return a.Ident() }
func (a *PrimFn) Equal(v Value) bool {
	a0, ok := v.(*PrimFn)
	if !ok {
		return false
	}
	return a0.String() == a.String()
}
func (a *PrimFn) Proto() *proto.Value {
	return &proto.Value{&proto.Value_PrimFn{&proto.PrimFn{a.Ident()}}}
}

type Ref uint64

func NewRef(u uint64) *Ref    { res := Ref(u); return &res }
func (*Ref) Type() ValueType  { return TypeRef }
func (r *Ref) Uint() uint64   { return uint64(*r) }
func (r *Ref) String() string { return "#" + strconv.FormatUint(r.Uint(), 10) }
func (r *Ref) Equal(v Value) bool {
	r0, ok := v.(*Ref)
	if !ok {
		return false
	}
	return r0.Uint() == r.Uint()
}
func (r *Ref) Proto() *proto.Value {
	return &proto.Value{&proto.Value_Ref{&proto.Ref{r.Uint()}}}
}

type State struct {
	Env  Env
	Refs *Intmap
}

func NewState(env Env, refs *Intmap) *State {
	return &State{
		Env:  env,
		Refs: refs,
	}
}
func (*State) Type() ValueType { return TypeState }
func (s *State) String() string {
	return "" // XXX
}
func (s *State) Equal(v Value) bool {
	return false // XXX
}
func (s *State) Proto() *proto.Value {
	/*
		return &proto.Value{&proto.Value_State{
			s.Env.Proto(),
			s.Refs.Proto(),
		}}
	*/
	return nil //XXX
}

func Unproto(pv *proto.Value) Value {
	if pv == nil {
		return nil
	}

	switch pv := pv.GetValue().(type) {
	case *proto.Value_Atom:
		return NewAtom(NewIdent(pv.Atom.Atom))
	case *proto.Value_Keyword:
		return NewKeyword(NewIdent(pv.Keyword.Keyword))
	case *proto.Value_String_:
		return NewString(NewIdent(pv.String_.String_))
	case *proto.Value_Boolean:
		return NewBool(pv.Boolean.Boolean)
	case *proto.Value_Num:
		return NewNum(pv.Num.Num)
	case *proto.Value_List:
		l := pv.List
		vs := make([]Value, len(l.Values))
		for i, v := range l.Values {
			vs[i] = Unproto(v)
		}
		return NewList(vs...)
	case *proto.Value_Vector:
		l := pv.Vector
		vs := make([]Value, len(l.Values))
		for i, v := range l.Values {
			vs[i] = Unproto(v)
		}
		return NewVector(vs...)
	case *proto.Value_PrimFn:
		res := PrimFn(pv.PrimFn.Fn)
		return &res
	case *proto.Value_Dict:
		d := pv.Dict
		m := Dict(make(map[Value]Value))
		for _, kvp := range d.Pairs {
			m[Unproto(kvp.Key)] = Unproto(kvp.Value)
		}
		return &m
	case *proto.Value_Ref:
		return NewRef(pv.Ref.Ref)
	case *proto.Value_Lambda:
		l := pv.Lambda
		bodies := make([]Value, len(l.Bodies))
		for i, body := range l.Bodies {
			bodies[i] = Unproto(body)
		}
		return NewLambda(l.Args, bodies,
			Unproto(&proto.Value{&proto.Value_Env{l.Env}}).(Env))
	case *proto.Value_LambdaRec:
		return NewLambdaRec(pv.LambdaRec.Self,
			Unproto(&proto.Value{&proto.Value_Lambda{pv.LambdaRec.Lambda}}).(*Lambda))
	case *proto.Value_Env:
		d := pv.Env
		m := NewListEnv()
		for _, kvp := range d.Pairs {
			m.Set(kvp.Key, Unproto(kvp.Value))
		}
		return m
	case *proto.Value_State:
		/*
			return NewState(
				Unproto(&proto.Value{&proto.Value_Env{pv.State.Env}}).(Env),
				Unproto(&proto.Value{&proto.Value_State{pv.State.State}}).(*State),
			)
		*/
		return nil //XXX
	default:
		return nil
	}
}
