package types

type Intmap struct {
	m    map[uint64]Value
	next uint64
}

func NewIntmap() *Intmap {
	return &Intmap{
		m:    make(map[uint64]Value),
		next: 0,
	}
}

func (m *Intmap) Insert(v Value) uint64 {
	ix := m.next
	m.m[ix] = v
	m.next = ix + 1
	return ix
}

func (m *Intmap) Get(ix uint64) (res Value, ok bool) {
	res, ok = m.m[ix]
	return
}

func (m *Intmap) Set(ix uint64, v Value) {
	m.m[ix] = v
}
