package pq

import (
	"container/heap"
	"golang.org/x/exp/constraints"
)

/*
 * This code is from https://golang.org/pkg/container/heap/
 */

type Number interface {
	constraints.Integer | constraints.Float
}

type Item[S any, T Number] struct {
	value    S // The value of the item; arbitrary.
	priority T // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[S any, T Number] []*Item[S, T]

func (pq PriorityQueue[S, T]) Len() int { return len(pq) }

func (pq PriorityQueue[S, T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue[S, T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[S, T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[S, T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[S, T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue[S, T]) update(item *Item[S, T], value S, priority T) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
