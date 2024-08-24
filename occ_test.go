package occ

import (
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {
	a := 123
	c := Wrap(&a)
	if *(c.Value()) != a {
		t.Fatalf("unexpected value %v, expected %v", *(c.Value()), a)
	}
	c.Update(func(old *int) *int {
		res := *old + 1
		return &res
	})
	if *(c.Value()) != a+1 {
		t.Fatalf("unexpected value %v, expected %v", *(c.Value()), a+1)
	}
}

func TestUpdateCollision(t *testing.T) {
	var childWG sync.WaitGroup
	childWG.Add(2)
	firstAcquiredPtr := make(chan struct{})
	secondAcquiredPtr := make(chan struct{})
	firstDone := make(chan struct{})
	var firstCalled, secondCalled int
	var firstOnce sync.Once
	var secondOnce sync.Once

	a := 0
	c := Wrap(&a)

	go func() {
		defer childWG.Done()
		defer close(firstDone)
		c.Update(func(old *int) *int {
			firstCalled++
			firstOnce.Do(func() { close(firstAcquiredPtr) })
			<-secondAcquiredPtr
			res := *old + 1
			return &res
		})
	}()

	go func() {
		defer childWG.Done()
		c.Update(func(old *int) *int {
			secondCalled++
			secondOnce.Do(func() { close(secondAcquiredPtr) })
			<-firstAcquiredPtr
			<-firstDone
			res := *old + 1
			return &res
		})
	}()

	childWG.Wait()
	if firstCalled != 1 {
		t.Fatalf("update txn in first goroutine called %d times, but should be called once", firstCalled)
	}
	if secondCalled != 2 {
		t.Fatalf("update txn in second goroutine called %d times, but should be called twice", secondCalled)
	}
	if *(c.Value()) != 2 {
		t.Fatalf("unexpected value: got %d instead of 2", 2)
	}
}
