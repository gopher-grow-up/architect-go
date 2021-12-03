package cirular_queue

import "errors"

type CircularQueue struct {
	Head  int
	Tail  int
	Items []string
	N     int
}

func NewCircularQueue(capacity int) CircularQueue {
	return CircularQueue{N: capacity, Items: make([]string, capacity)}
}

func (q CircularQueue) Enqueue(item string) error {
	if q.IsFull() {
		return errors.New("queue is full")
	}
	q.Items[q.Tail] = item
	q.Tail = (q.Tail + 1) % q.N
	return nil
}

func (q CircularQueue) IsFull() bool {
	return (q.Tail+1)%q.N == q.Head
}

func (q CircularQueue) Dequeue() string {
	if q.IsEmpty() {
		return ""
	}
	item := q.Items[q.Head]
	q.Head = (q.Head + 1) % q.N
	return item
}

func (q CircularQueue) IsEmpty() bool {
	return q.Head == q.Tail
}
