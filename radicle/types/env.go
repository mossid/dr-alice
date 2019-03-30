package types

type Env interface {
	Set(name Ident, value Value) Env
	Get(name Ident) (Value, bool)
	CloneMutable() Env
	CloneImmutable() Env
	IsImmutable() bool
}

type node struct {
	name  Ident
	value Value
	next  *node
}

type listEnv struct {
	top *node
}

var _ Env = (*listEnv)(nil)

func (env *listEnv) Set(name Ident, value Value) Env {
	res := &listEnv{}
	res.top = &node{name, value, env.top}
	return res
}

func (env *listEnv) Get(name Ident) (Value, bool) {
	ptr := env.top
	for ptr != nil {
		if ptr.name == name {
			return ptr.value, true
		}
		ptr = ptr.next
	}
	return nil, false
}

func (env *listEnv) CloneImmutable() Env {
	return env
}

func (env *listEnv) CloneMutable() Env {
	res := &mapEnv{make(map[Ident]Value)}
	ptr := env.top
	for ptr != nil {
		res.m[ptr.name] = ptr.value
		ptr = ptr.next
	}
	return res
}

func (env *listEnv) IsImmutable() bool {
	return true
}

type mapEnv struct {
	m map[Ident]Value
}

var _ Env = (*mapEnv)(nil)

func (env *mapEnv) Set(name Ident, value Value) Env {
	env.m[name] = value
	return env
}

func (env *mapEnv) Get(name Ident) (res Value, ok bool) {
	res, ok = env.m[name]
	return
}

func (env *mapEnv) CloneImmutable() Env {
	res := &listEnv{}
	for k, v := range env.m {
		res.top = &node{k, v, res.top}
	}
	return res
}

func (env *mapEnv) CloneMutable() Env {
	res := &mapEnv{make(map[Ident]Value)}
	for k, v := range env.m {
		res.m[k] = v
	}
	return res
}

func (env *mapEnv) IsImmutable() bool {
	return false
}
