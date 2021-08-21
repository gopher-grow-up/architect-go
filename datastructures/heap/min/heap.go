// Package min
// File     heap
// Created by lt on 2021/6/30
// Copyright © 2020-2020 lt. All rights reserved.

package min

type Heap struct {
	Item []int
	count int
}

func NewMinHeap(capacity int)*Heap  {
	return &Heap{Item: make([]int,capacity+1)}
}

func (h *Heap) Size()int {
	return h.count
}

func (h *Heap) Push(value int)  {
	h.Item[h.count+1] = value
	h.count +=1
	h.shiftUp(h.count)
}

func (h *Heap) Pop()int  {
	if h.count <= 0 {
		return 0
	}

	v := h.Item[1]
	h.Item[1],h.Item[h.count] = h.Item[h.count],h.Item[1]
	h.count -= 1

	h.shiftDown(1)
	return v

}

func (h *Heap) shiftUp(k int)  {

	for k > 1 && h.Item[k/2] > h.Item[k] {
		h.Item[k/2],h.Item[k] = h.Item[k],h.Item[k/2]
		k /= 2
	}
}

func (h *Heap) shiftDown(k int)  {
	for 2*k <=h.count {
		j := 2*k
		if j+1 <= h.count && h.Item[j+1] < h.Item[j] {
			j++
		}

		if h.Item[k] <= h.Item[j] {
			break
		}
		h.Item[k],h.Item[j] = h.Item[j],h.Item[k]
		k = j
	}
}
