// Package single
// File     single
// Created by lt on 2021/6/29
// Copyright © 2020-2020 lt. All rights reserved.

package single

type Node struct {
	Prev *Node
	Next *Node
	Val  int
}

type LinkedList struct {
	Head *Node
	Size int
}

// NewSingleLinkedList  Initialize your data structure here.
func NewSingleLinkedList() LinkedList {
	return LinkedList{
		Head: &Node{},
	}
}

// Get the value of the index-th node in the linked list. If the index is invalid, return -1.
func (l *LinkedList) Get(index int) int {
	if index >= l.Size || index < 0 {
		return -1
	}

	cur := l.Head
	for i := 0; i < index+1; i++ {
		cur = cur.Next
	}

	return cur.Val
}

// AddAtHead  Add a node of value val before the first element of the linked list. After the insertion, the new node will be the first node of the linked list.
func (l *LinkedList) AddAtHead(val int) {
	l.AddAtIndex(0, val)
}

// AddAtTail  Append a node of value val to the last element of the linked list.
func (l *LinkedList) AddAtTail(val int) {
	l.AddAtIndex(l.Size, val)
}

// AddAtIndex  Add a node of value val before the index-th node in the linked list. If index equals to the length of linked list, the node will be appended to the end of linked list. If index is greater than the length, the node will not be inserted.
func (l *LinkedList) AddAtIndex(index int, val int) {
	if index > l.Size || index < 0 {
		return
	}
	l.Size++
	cur := l.Head

	for i := 0; i < index; i++ {
		cur = cur.Next
	}
	node := &Node{
		Prev: cur,
		Next: cur.Next,
		Val:  val,
	}
	cur.Next = node
}

// DeleteAtIndex  Delete the index-th node in the linked list, if the index is valid.
func (l *LinkedList) DeleteAtIndex(index int) {
	if index >= l.Size || index < 0 {
		return
	}
	l.Size--
	cur := l.Head

	for i := 0; i < index; i++ {
		cur = cur.Next
	}

	cur.Next = cur.Next.Next
	if cur.Next != nil && cur.Next.Next != nil {
		cur.Next.Next.Prev = cur
	}
}
