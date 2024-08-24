package occ

import "sync/atomic"

type Container[T any] atomic.Pointer[T]

func Wrap[T any](ptr *T) *Container[T] {
	c := new(atomic.Pointer[T])
	c.Store(ptr)
	return (*Container[T])(c)
}

func (c *Container[T]) Value() *T {
	return (*atomic.Pointer[T])(c).Load()
}

func (c *Container[T]) Update(txn func(oldValue *T) *T) {
	for {
		oldValue := c.Value()
		newValue := txn(oldValue)
		if (*atomic.Pointer[T])(c).CompareAndSwap(oldValue, newValue) {
			break
		}
	}
}
