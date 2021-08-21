// Package single
// File     single_test.go
//
// Created by lt on 2021/6/29
// Copyright © 2020-2020 lt. All rights reserved.

package single

import (
	"reflect"
	"testing"
)

func TestConstructor(t *testing.T) {
	tests := []struct {
		name string
		want LinkedList
	}{
		{
			want: NewSingleLinkedList(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSingleLinkedList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSingleLinkedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkedList_AddAtHead(t *testing.T) {
	type args struct {
		val int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{val: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewSingleLinkedList()
			l.AddAtHead(tt.args.val)
			if got := l.Get(0); got != tt.args.val {
				t.Errorf("AddAtHead() = %v, want %v", got, tt.args.val)
			}
		})
	}
}

func TestLinkedList_AddAtIndex(t *testing.T) {
	type args struct {
		index int
		val   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				index: 3,
				val:   4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSingleLinkedList()
			for i := 0; i < 10; i++ {
				s.AddAtIndex(i, i)
			}
			s.AddAtIndex(tt.args.index, tt.args.val)
			if got := s.Get(tt.args.index); got != tt.args.val {
				t.Errorf("AddAtIndex() = %v, want %v", got, tt.args.val)
			}
		})
	}
}

func TestLinkedList_AddAtTail(t *testing.T) {
	type args struct {
		val int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{val: 33},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSingleLinkedList()
			for i := 0; i < 10; i++ {
				s.AddAtIndex(i, i)
			}
			s.AddAtTail(tt.args.val)
			if got := s.Get(10); got != tt.args.val {
				t.Errorf("AddAtIndex() = %v, want %v", got, tt.args.val)
			}
		})
	}
}

func TestLinkedList_DeleteAtIndex(t *testing.T) {

	type args struct {
		index int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{index: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSingleLinkedList()
			for i := 0; i < 10; i++ {
				s.AddAtIndex(i, i)
			}

			s.DeleteAtIndex(1)
			if got := s.Get(1); got != 2 {
				t.Errorf("DeleteAtIndex() = %v, want %v", got, 2)
			}

		})
	}
}

func TestLinkedList_Get(t *testing.T) {
	type fields struct {
		Head *Node
		Size int
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			args: args{index: 2},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSingleLinkedList()
			for i := 0; i < 10; i++ {
				s.AddAtIndex(i, i)
			}
			if got := s.Get(tt.args.index); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
