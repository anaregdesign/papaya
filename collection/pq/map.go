package pq

import (
	"container/heap"
)

func FilterMap[S comparable, T Number](m map[S]T, top int) map[S]T {
	if len(m) <= top {
		return m
	}

	pq := make(PriorityQueue[S, T], len(m))

	i := 0
	for k, v := range m {
		pq[i] = &Item[S, T]{
			value:    k,
			priority: v,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	filtered := make(map[S]T)
	for i := 0; i < top; i++ {
		item := heap.Pop(&pq).(*Item[S, T])
		filtered[item.value] = item.priority
	}
	return filtered
}
