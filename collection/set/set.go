package set

import (
	"github.com/anaregdesign/papaya/model/function"
)

type Set[T comparable] struct {
	set map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		set: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(value T) {
	s.set[value] = struct{}{}
}

func (s *Set[T]) Remove(value T) {
	delete(s.set, value)
}

func (s *Set[T]) Has(value T) bool {
	_, ok := s.set[value]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.set)
}

func (s *Set[T]) Clear() {
	s.set = make(map[T]struct{})
}

func (s *Set[T]) Values() []T {
	values := make([]T, 0, len(s.set))
	for k, _ := range s.set {
		values = append(values, k)
	}
	return values
}

func (s *Set[T]) ForEach(consumer function.Consumer[T]) {
	for value := range s.set {
		consumer(value)
	}
}

func (s *Set[T]) Filter(predicate function.Predicate[T]) {
	for k, _ := range s.set {
		if !predicate(k) {
			delete(s.set, k)
		}
	}
}

func (s *Set[T]) Reduce(operator function.Operator[T]) T {
	var result T
	for value := range s.set {
		result = operator(result, value)
	}
	return result
}

func (s *Set[T]) AnyMatch(predicate function.Predicate[T]) bool {
	for value := range s.set {
		if predicate(value) {
			return true
		}
	}
	return false
}

func (s *Set[T]) AllMatch(predicate function.Predicate[T]) bool {
	for value := range s.set {
		if !predicate(value) {
			return false
		}
	}
	return true
}

func (s *Set[T]) NoneMatch(predicate function.Predicate[T]) bool {
	for value := range s.set {
		if predicate(value) {
			return false
		}
	}
	return true
}
