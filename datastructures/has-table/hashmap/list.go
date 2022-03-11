package hashmap

import (
	"sync/atomic"
	"unsafe"
)

type List struct {
	count uintptr
	head  *ListElement
}

func NewList() *List {
	return &List{head: &ListElement{}}
}

func (l *List) Len() int {
	if l == nil {
		return 0
	}

	return int(atomic.LoadUintptr(&l.count))
}

func (l *List) Head() *ListElement {
	if l == nil {
		return nil
	}

	return l.head
}

func (l *List) First() *ListElement {
	if l == nil {
		return nil
	}

	return l.head.Next()
}

func (l *List) Add(element *ListElement, searchStart *ListElement) (existed bool, inserted bool) {
	left, found, right := l.search(searchStart, element)
	if found != nil {
		return true, false
	}

	return false, l.insertAt(element, left, right)
}

func (l *List) AddOrUpdate(element *ListElement, searchStart *ListElement) bool {
	left, found, right := l.search(searchStart, element)
	if found != nil {
		found.setValue(element.value)
		return true
	}

	return l.insertAt(element, left, right)
}

func (l *List) Cas(element *ListElement, oldValue interface{}, searchStart *ListElement) bool {
	_, found, _ := l.search(searchStart, element)
	if found == nil {
		return false
	}

	if found.casValue(oldValue, element.value) {
		return true
	}

	return false
}

func (l *List) Delete(element *ListElement) {
	if !atomic.CompareAndSwapUintptr(&element.deleted, uintptr(0), uintptr(1)) {
		return
	}

	for {
		left := element.Previous()
		right := element.Next()

		if left == nil {
			if !atomic.CompareAndSwapPointer(&l.head.nextElement, unsafe.Pointer(element), unsafe.Pointer(right)) {
				continue
			}
		} else {
			if !atomic.CompareAndSwapPointer(&left.nextElement, unsafe.Pointer(element), unsafe.Pointer(right)) {
				continue
			}
		}

		if right != nil {
			atomic.CompareAndSwapPointer(&right.previousElement, unsafe.Pointer(element), unsafe.Pointer(left))
		}
		break
	}
}

func (l *List) insertAt(element *ListElement, left *ListElement, right *ListElement) bool {
	if left == nil {
		element.previousElement = unsafe.Pointer(l.head)
		element.nextElement = unsafe.Pointer(right)
		if !atomic.CompareAndSwapPointer(&l.head.nextElement, unsafe.Pointer(right), unsafe.Pointer(element)) {
			return false
		}

		if right != nil {
			if !atomic.CompareAndSwapPointer(&right.previousElement, unsafe.Pointer(l.head), unsafe.Pointer(element)) {
				return false
			}
		}
	} else {
		element.previousElement = unsafe.Pointer(left)
		element.nextElement = unsafe.Pointer(right)

		if !atomic.CompareAndSwapPointer(&left.nextElement, unsafe.Pointer(right), unsafe.Pointer(element)) {
			return false
		}

		if right != nil {
			if !atomic.CompareAndSwapPointer(&right.previousElement, unsafe.Pointer(left), unsafe.Pointer(element)) {
				return false
			}
		}
	}
	atomic.AddUintptr(&l.count, 1)
	return true
}

func (l *List) search(searchStart *ListElement, item *ListElement) (left *ListElement, found *ListElement, right *ListElement) {
	if searchStart != nil && item.keyHash < searchStart.keyHash {
		searchStart = nil
	}
	if searchStart == nil {
		left = l.head
		found = left.Next()
		if found == nil {
			return nil, nil, nil
		}
	} else {
		found = searchStart
	}

	for {
		if item.keyHash == found.keyHash {
			return nil, found, nil
		}

		if item.keyHash < found.keyHash {
			if l.head == left {
				return nil, nil, found
			}

			return left, nil, found
		}

		left = found
		found = left.Next()
		if found == nil {
			return left, nil, nil
		}
	}
}
