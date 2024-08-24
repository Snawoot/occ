package occ_test

import (
	"fmt"
	"sync"

	"github.com/Snawoot/occ"
)

func Example() {
	ctr := occ.Wrap(new(int))
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctr.Update(func(old *int) *int {
				res := *old + 1
				return &res
			})
		}()
	}
	fmt.Printf("ctr = %d\n", *(ctr.Value()))
}
