package gobase

import "sync/atomic"

type Int64 struct {
	v int64
}

func (a *Int64) Int64() int64 {
	return atomic.LoadInt64(&a.v)
}

func (a *Int64) AsInt() int {
	return int(a.Int64())
}

func (a *Int64) Set(v int64) {
	atomic.StoreInt64(&a.v, v)
}

func (a *Int64) CompareAndSwap(o, n int64) bool {
	return atomic.CompareAndSwapInt64(&a.v, o, n)
}

func (a *Int64) Swap(v int64) int64 {
	return atomic.SwapInt64(&a.v, v)
}

func (a *Int64) Add(v int64) int64 {
	return atomic.AddInt64(&a.v, v)
}

func (a *Int64) Sub(v int64) int64 {
	return a.Add(-v)
}

func (a *Int64) Incr() int64 {
	return a.Add(1)
}

func (a *Int64) Decr() int64 {
	return a.Add(-1)
}

type Bool struct {
	c Int64
}

func (b *Bool) Bool() bool {
	return b.IsTrue()
}

func (b *Bool) IsTrue() bool {
	return b.c.Int64() != 0
}

func (b *Bool) IsFalse() bool {
	return b.c.Int64() == 0
}

func (b *Bool) toInt64(v bool) int64 {
	if v {
		return 1
	} else {
		return 0
	}
}

func (b *Bool) Set(v bool) {
	b.c.Set(b.toInt64(v))
}

func (b *Bool) CompareAndSwap(o, n bool) bool {
	return b.c.CompareAndSwap(b.toInt64(o), b.toInt64(n))
}

func (b *Bool) Swap(v bool) bool {
	return b.c.Swap(b.toInt64(v)) != 0
}

