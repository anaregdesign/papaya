package graph

import "time"

type weightValue struct {
	value float64
	ttl   time.Time
}

func (w weightValue) expired() bool {
	return time.Now().After(w.ttl)
}

type weight struct {
	values []weightValue
}

func newWeight() *weight {
	return &weight{
		values: make([]weightValue, 0),
	}
}

func (w *weight) value() float64 {
	w.flush()
	var sum float64
	for _, v := range w.values {
		sum += v.value
	}
	return sum
}

func (w *weight) add(value float64, ttl time.Duration) {
	w.values = append(w.values, weightValue{
		value: value,
		ttl:   time.Now().Add(ttl),
	})
}

func (w *weight) isZero() bool {
	return w.value() == 0
}

func (w *weight) flush() {
	v := make([]weightValue, 0)
	for _, value := range w.values {
		if !value.expired() {
			v = append(v, value)
		}
	}
	w.values = v
}
