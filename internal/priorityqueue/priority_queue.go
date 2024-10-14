package priorityqueue

import (
	"container/heap"
	"orderbook/internal/models"
)

type PriorityQueue []*models.Order

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Price > pq[j].Price
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	heap.Push(pq, x.(*models.Order))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityQueue) Peek() any {
	return pq[0]
}

func InitPriorityQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &pq
}
