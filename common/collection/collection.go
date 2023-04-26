package collection

import (
	"github.com/anaregdesign/papaya/common/concurrent"
	"github.com/anaregdesign/papaya/common/model"
	"runtime"
)

func ForEach[T any](slice []T, consumer model.Consumer[T]) {
	s := concurrent.NewAwaitableSemaphore(runtime.NumCPU())

	for _, element := range slice {
		s.Acquire()
		go func() {
			defer s.Release()
			consumer(element)
		}()
	}

	s.Wait()
}
