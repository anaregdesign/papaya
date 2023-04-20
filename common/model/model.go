package model

type Consumer[T any] func(T)
type Supplier[T any] func() T
type Predicate[T any] func(T) bool
type Function[T, R any] func(T) R

type Loader[S, T any] func(S) (T, bool)
