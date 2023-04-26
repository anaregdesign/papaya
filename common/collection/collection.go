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

func Map[S, T any](slice []S, function model.Function[S, T]) []T {
	s := concurrent.NewAwaitableSemaphore(runtime.NumCPU())
	result := make([]T, len(slice))

	for i, element := range slice {
		s.Acquire()
		go func(i int, element S) {
			defer s.Release()
			result[i] = function(element)
		}(i, element)
	}

	s.Wait()
	return result
}
