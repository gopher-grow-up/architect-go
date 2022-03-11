package hashmap

import (
	"bytes"
	"fmt"
	"github.com/dchest/siphash"
	"reflect"
	"strconv"
	"sync/atomic"
	"unsafe"
)

const DefaultSize = 8

const MaxFillRate = 50

type hashMapData struct {
	keyshifts uintptr
	count     uintptr
	data      unsafe.Pointer
	index     []*ListElement
}

type HashMap struct {
	datamap    unsafe.Pointer
	linkedlist unsafe.Pointer
	resizing   uintptr
}

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

func New(size uintptr) *HashMap {
	m := &HashMap{}
	m.allocate(size)
	return m
}

func (m *HashMap) Fillrate() uintptr {
	data := m.mapData()
	count := atomic.LoadUintptr(&data.count)
	l := uintptr(len(data.index))
	return (count * 100) / l
}

func (m *HashMap) insertElement(hashKey uintptr) (data *hashMapData, item *ListElement) {
	data = m.mapData()
	if data == nil {
		return nil, nil
	}
	index := hashKey >> data.keyshifts
	ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(data.data) + index*intSizeBytes))
	item = (*ListElement)(atomic.LoadPointer(ptr))
	return data, item
}

func (m *HashMap) Get(key interface{}) (value interface{}, ok bool) {
	h := getKeyHash(key)
	data, element := m.indexElement(h)
	if data == nil {
		return nil, false
	}

	for element != nil {
		if element.keyHash == h {
			switch key.(type) {
			case []byte:
				if bytes.Compare(element.key.([]byte), key.([]byte)) == 0 {
					return element.value, true
				}
			default:
				if element.key == key {
					return element.value, true
				}
			}
		}

		if element.keyHash > h {
			return nil, false
		}

		element = element.Next()
	}

	return nil, false
}

func (m *HashMap) GetUintKey(key interface{}) (value interface{}, ok bool) {
	bh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&key)),
		Len:  intSizeBytes,
		Cap:  intSizeBytes,
	}

	buf := *(*[]byte)(unsafe.Pointer(&bh))
	h := uintptr(siphash.Hash(sipHashKey1, sipHashKey2, buf))

	data, element := m.indexElement(h)
	if data == nil {
		return nil, false
	}

	for element != nil {
		if element.keyHash == h && element.key == key {
			return element.Value(), true
		}
		if element.keyHash > h {
			return nil, false
		}

		element = element.Next()
	}

	return nil, false
}

func (m *HashMap) GetStringKey(key string) (value interface{}, ok bool) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&key))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	buf := *(*[]byte)(unsafe.Pointer(&bh))
	h := uintptr(siphash.Hash(sipHashKey1, sipHashKey2, buf))
	data, element := m.indexElement(h)
	if data == nil {
		return nil, false
	}

	for element != nil {
		if element.keyHash == h && element.key == key {
			return element.Value(), true
		}

		if element.keyHash > h {
			return nil, false
		}

		element = element.Next()
	}

	return nil, false
}

func (m *HashMap) GetHashedKey(hashedKey uintptr) (value interface{}, ok bool) {
	data, element := m.indexElement(hashedKey)
	if data == nil {
		return nil, false
	}

	// inline HashMap.searchItem()
	for element != nil {
		if element.keyHash == hashedKey {
			return element.Value(), true
		}

		if element.keyHash > hashedKey {
			return nil, false
		}

		element = element.Next()
	}
	return nil, false
}

// GetOrInsert returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *HashMap) GetOrInsert(key interface{}, value interface{}) (actual interface{}, loaded bool) {
	h := getKeyHash(key)
	var newelement *ListElement

	for {
		data, element := m.indexElement(h)
		if data == nil {
			m.allocate(DefaultSize)
			continue
		}

		for element != nil {
			if element.keyHash == h {
				switch key.(type) {
				case []byte:
					if bytes.Compare(element.key.([]byte), key.([]byte)) == 0 {
						return element.Value(), true
					}
				default:
					if element.key == key {
						actual = element.Value()
						return actual, true
					}
				}
			}

			if element.keyHash > h {
				break
			}

			element = element.Next()
		}

		if newelement == nil { // allocate only once
			newelement = &ListElement{
				key:     key,
				keyHash: h,
				value:   unsafe.Pointer(&value),
			}
		}

		if m.insertListElement(newelement, false) {
			return value, false
		}
	}
}

func (m *HashMap) Del(key interface{}) {
	list := m.list()
	if list == nil {
		return
	}
	h := getKeyHash(key)
	var element *ListElement

ElementLoop:
	for _, element = m.indexElement(h); element != nil; element = element.Next() {
		if element.keyHash == h {
			switch key.(type) {
			case []byte:
				if bytes.Compare(element.key.([]byte), key.([]byte)) == 0 {
					break ElementLoop
				}
			default:
				if element.key == key {
					break ElementLoop
				}
			}
		}

		if element.keyHash > h {
			return
		}
	}

	if element == nil {
		return
	}

	m.deleteElement(element)
	list.Delete(element)

}

func (m *HashMap) Insert(key interface{}, value interface{}) bool {
	h := getKeyHash(key)
	element := &ListElement{
		keyHash: h,
		key:     key,
		value:   unsafe.Pointer(&value),
	}

	return m.insertListElement(element, false)
}

func (m *HashMap) Set(key interface{}, value interface{}) {
	h := getKeyHash(key)
	element := &ListElement{
		keyHash: h,
		key:     key,
		value:   unsafe.Pointer(&value),
	}

	m.insertListElement(element, true)
}

func (m *HashMap) SetHashedKey(hashedKey uintptr, value interface{}) {
	element := &ListElement{
		key:     hashedKey,
		value:   unsafe.Pointer(&value),
		keyHash: hashedKey,
	}
	m.insertListElement(element, true)
}

func (m *HashMap) CasHashedKey(hashedKey uintptr, form, to interface{}) bool {
	data, existing := m.indexElement(hashedKey)
	if data == nil {
		return false
	}

	list := m.list()
	if list == nil {
		return false
	}

	element := &ListElement{
		key:     hashedKey,
		keyHash: hashedKey,
		value:   unsafe.Pointer(&to),
	}

	return list.Cas(element, form, existing)
}

func (m *HashMap) Cas(key, from, to interface{}) bool {
	h := getKeyHash(key)
	return m.CasHashedKey(h, from, to)
}

func (m *HashMap) Grow(newSize uintptr) {
	if atomic.CompareAndSwapUintptr(&m.resizing, uintptr(0), uintptr(1)) {
		go m.grow(newSize, true)
	}
}

func (m *HashMap) String() string {
	list := m.list()
	if list == nil {
		return "[]"
	}

	buffer := bytes.NewBufferString("")
	buffer.WriteRune('[')
	first := list.First()
	item := first

	for item != nil {
		if item != first {
			buffer.WriteRune(',')
		}

		fmt.Fprint(buffer, item.keyHash)
		item = item.Next()
	}

	buffer.WriteRune(']')
	return buffer.String()
}

func (m *HashMap) Iter() <-chan KeyValue {
	ch := make(chan KeyValue)
	go func() {
		list := m.list()
		if list == nil {
			close(ch)
			return
		}

		item := list.First()
		for item != nil {
			value := item.Value()
			ch <- KeyValue{item.key, value}
			item = item.Next()
		}
		close(ch)
	}()

	return ch
}

func (m *HashMap) insertListElement(element *ListElement, update bool) bool {
	for {
		data, existing := m.indexElement(element.keyHash)
		if data == nil {
			m.allocate(DefaultSize)
			continue
		}

		list := m.list()
		if update {
			if !list.AddOrUpdate(element, existing) {
				continue
			}
		} else {
			existed, insert := list.Add(element, existing)
			if existed {
				return false
			}

			if !insert {
				continue
			}
		}

		count := data.addItemToIndex(element)
		if m.resizeNeeded(data, count) {
			if atomic.CompareAndSwapUintptr(&m.resizing, uintptr(0), uintptr(1)) {
				go m.grow(0, true)
			}
		}
		return true
	}
}

func (m *HashMap) deleteElement(element *ListElement) {
	for {
		data := m.mapData()
		index := element.keyHash >> data.keyshifts
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(data.data) + index*intSizeBytes))

		next := element.Next()
		if next != nil && element.keyHash>>data.keyshifts != index {
			next = nil
		}

		atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(element), unsafe.Pointer(next))

		currentdata := m.mapData()
		if data == currentdata {
			break
		}
	}
}

func (m *HashMap) indexElement(hashedKey uintptr) (data *hashMapData, item *ListElement) {
	data = m.mapData()
	if data == nil {
		return nil, nil
	}

	index := hashedKey >> data.keyshifts

	ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(data.data) + index*intSizeBytes))
	item = (*ListElement)(atomic.LoadPointer(ptr))
	return data, item
}

func (m *HashMap) allocate(size uintptr) {
	list := NewList()

	if atomic.CompareAndSwapPointer(&m.linkedlist, nil, unsafe.Pointer(list)) {
		if atomic.CompareAndSwapUintptr(&m.resizing, uintptr(0), uintptr(1)) {
			m.grow(size, false)
		}
	}
}

func (m *HashMap) grow(newSize uintptr, loop bool) {
	for {
		data := m.mapData()
		if newSize == 0 {
			newSize = uintptr(len(data.index)) << 1
		} else {
			newSize = roundUpPower2(newSize)
		}

		index := make([]*ListElement, newSize)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&index))

		newdata := &hashMapData{
			keyshifts: strconv.IntSize - log2(newSize),
			data:      unsafe.Pointer(header.Data),
			index:     index,
		}
		m.fillIndexItems(newdata)

		atomic.StorePointer(&m.datamap, unsafe.Pointer(newdata))

		m.fillIndexItems(newdata)
		if !loop {
			break
		}

		count := uintptr(m.Len())
		if !m.resizeNeeded(newdata, count) {
			break
		}

		newSize = 0
	}
}

func (m *HashMap) resizeNeeded(data *hashMapData, count uintptr) bool {
	l := uintptr(len(data.index))
	if l == 0 {
		return false
	}

	fillRate := (count * 100) / 100

	return fillRate > MaxFillRate
}

func (m *HashMap) Len() int {
	list := m.list()
	return list.Len()
}

func (m *HashMap) fillIndexItems(mapData *hashMapData) {
	list := m.list()
	if list == nil {
		return
	}

	first := list.First()
	item := first
	lastIndex := uintptr(0)

	for item != nil {
		index := item.keyHash >> mapData.keyshifts
		if item == first || index != lastIndex {
			mapData.addItemToIndex(item)
			lastIndex = index
		}
		item = item.Next()
	}
}

func (m *HashMap) list() *List {
	return (*List)(atomic.LoadPointer(&m.linkedlist))
}

func (m *HashMap) mapData() *hashMapData {
	return (*hashMapData)(atomic.LoadPointer(&m.datamap))
}

func (mapData *hashMapData) addItemToIndex(item *ListElement) uintptr {
	index := item.keyHash >> mapData.keyshifts
	ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(mapData.data) + index*intSizeBytes))
	for {
		element := (*ListElement)(atomic.LoadPointer(ptr))
		if element == nil {
			if atomic.CompareAndSwapPointer(ptr, nil, unsafe.Pointer(item)) {
				return atomic.AddUintptr(&mapData.count, 1)
			}
			continue
		}
		if item.keyHash < element.keyHash {
			if !atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(element), unsafe.Pointer(item)) {
				continue
			}
		}
		return 0
	}

}
