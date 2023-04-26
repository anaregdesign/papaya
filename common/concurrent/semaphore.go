package concurrent

import "sync"

type AwaitableSemaphore struct {
	permits chan struct{}
	wg      sync.WaitGroup
}

func NewAwaitableSemaphore(permits int) *AwaitableSemaphore {
	return &AwaitableSemaphore{
		permits: make(chan struct{}, permits),
	}
}

func (s *AwaitableSemaphore) Acquire() {
	s.wg.Add(1)
	s.permits <- struct{}{}
}

func (s *AwaitableSemaphore) Release() {
	<-s.permits
	s.wg.Done()
}

func (s *AwaitableSemaphore) Wait() {
	s.wg.Wait()
}
