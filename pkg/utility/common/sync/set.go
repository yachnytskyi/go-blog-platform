package sync

import (
	"sync"
)

// Set represents a synchronized set implemented using a map.
type Set[T comparable] struct {
	mutex sync.RWMutex
	data  map[T]struct{}
}

func NewSet[T comparable](capacity uint) *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}, capacity),
	}
}

func (set *Set[T]) Add(value T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	set.data[value] = struct{}{}
}

func (set *Set[T]) AddAll(values ...T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	for _, value := range values {
		set.data[value] = struct{}{}
	}
}

func (set *Set[T]) Contains(value T) bool {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	_, exists := set.data[value]
	return exists
}

func (set *Set[T]) IsUnique(value T) bool {
	_, exists := set.data[value]
	return !exists
}

func (set *Set[T]) Len() int {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	return len(set.data)
}

func (set *Set[T]) Values() []T {
	values := make([]T, 0, set.Len())
	for value := range set.data {
		values = append(values, value)
	}

	return values
}

func (set *Set[T]) Delete(value T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	delete(set.data, value)
}

func (set *Set[T]) Clear() {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	set.data = make(map[T]struct{}, set.Len())
}

// OrderedSet represents an ordered set implemented using both a map and a slice.
type OrderedSet[T comparable] struct {
	mutex    sync.RWMutex
	setMap   map[T]struct{}
	setSlice []T
}

func NewOrderedSet[T comparable](capacity int) *OrderedSet[T] {
	return &OrderedSet[T]{
		setMap:   make(map[T]struct{}, capacity),
		setSlice: make([]T, 0, capacity),
	}
}

func (orderedSet *OrderedSet[T]) Add(value T) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()
	if orderedSet.IsUnique(value) {
		orderedSet.setMap[value] = struct{}{}
		orderedSet.setSlice = append(orderedSet.setSlice, value)
	}
}

func (orderedSet *OrderedSet[T]) AddAll(values ...T) {
	for _, value := range values {
		orderedSet.Add(value)
	}
}

func (orderedSet *OrderedSet[T]) Contains(value T) bool {
	orderedSet.mutex.RLock()
	defer orderedSet.mutex.RUnlock()
	_, exists := orderedSet.setMap[value]
	return exists
}

func (orderedSet *OrderedSet[T]) IsUnique(value T) bool {
	orderedSet.mutex.RLock()
	defer orderedSet.mutex.RUnlock()
	_, exists := orderedSet.setMap[value]
	return !exists
}

func (orderedSet *OrderedSet[T]) Values() []T {
	orderedSet.mutex.RLock()
	defer orderedSet.mutex.RUnlock()
	values := make([]T, len(orderedSet.setSlice))
	copy(values, orderedSet.setSlice)
	return values
}

func (orderedSet *OrderedSet[T]) Delete(index uint) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()
	if int(index) < len(orderedSet.setSlice) {
		value := orderedSet.setSlice[index]
		delete(orderedSet.setMap, value)
		orderedSet.setSlice = append(orderedSet.setSlice[:index], orderedSet.setSlice[index+1:]...)
	}
}

func (orderedSet *OrderedSet[T]) Clear() {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()

	orderedSet.setMap = make(map[T]struct{}, len(orderedSet.setMap))
	orderedSet.setSlice = orderedSet.setSlice[:0]
}
