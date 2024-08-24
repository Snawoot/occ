// Package occ provides generic optimistic concurrency wrapper for values which
// can be represented by a pointer and accessed or replaced with a new copy
// concurrently.
//
// It is useful for implementation of lock-free concurrency on top of immutable
// data structures.
package occ

import "sync/atomic"

// Container uses optimistic concurrency control for access and update of
// pointer it wraps. Essentially it is just an atomic pointer with
// compare-and-set logic around updates, so concurrent writes will no get
// lost and colliding changes will be retried instead.
type Container[T any] atomic.Pointer[T]

// Wrap creates Container from a pointer to any type. This pointer should not
// be used directly after wrapping.
func Wrap[T any](ptr *T) *Container[T] {
	c := new(atomic.Pointer[T])
	c.Store(ptr)
	return (*Container[T])(c)
}

// Value returns snapshot of wrapped value.
func (c *Container[T]) Value() *T {
	return (*atomic.Pointer[T])(c).Load()
}

// Update runs a transaction function txn which accepts old pointer value and
// returns new pointer value. txn function must not modify value referenced by
// an old pointer, instead it should return a new pointer with modified copy of
// old value.
func (c *Container[T]) Update(txn func(oldValue *T) (newValue *T)) {
	for {
		oldValue := c.Value()
		newValue := txn(oldValue)
		if (*atomic.Pointer[T])(c).CompareAndSwap(oldValue, newValue) {
			break
		}
	}
}
