package sync

import (
	"sync"
)

// Set represents a synchronized set implemented using a map.
type Set[T comparable] struct {
	mutex sync.RWMutex
	data  map[T]struct{}
}

// NewSet creates a new set with an optional capacity.
//
// Parameters:
// - capacity: The initial capacity of the set.
//
// Returns:
// - A pointer to the newly created Set.
func NewSet[T comparable](capacity uint) *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}, capacity),
	}
}

// Add adds a value to the synchronized set.
//
// Parameters:
// - value: The value to be added to the set.
func (set *Set[T]) Add(value T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	set.data[value] = struct{}{}
}

// AddAll adds multiple values to the synchronized set.
//
// Parameters:
// - values: The values to be added to the set.
func (set *Set[T]) AddAll(values ...T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	for _, value := range values {
		set.data[value] = struct{}{}
	}
}

// Contains checks if a value exists in the set.
//
// Parameters:
// - value: The value to check for existence in the set.
//
// Returns:
// - A boolean indicating whether the value is present in the set.
func (set *Set[T]) Contains(value T) bool {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	_, exists := set.data[value]
	return exists
}

// IsUnique checks if a value is unique in the set.
//
// Parameters:
// - value: The value to check for uniqueness in the set.
//
// Returns:
// - A boolean indicating whether the value is unique in the set.
func (set *Set[T]) IsUnique(value T) bool {
	_, exists := set.data[value]
	return !exists
}

// Len returns the number of values in the set.
//
// Returns:
// - The number of values in the set.
func (set *Set[T]) Len() int {
	set.mutex.RLock()
	defer set.mutex.RUnlock()
	return len(set.data)
}

// Values returns the values from the synchronized set.
//
// Returns:
// - A slice containing all the values in the set.
func (set *Set[T]) Values() []T {
	values := make([]T, 0, set.Len())
	for value := range set.data {
		values = append(values, value)
	}

	return values
}

// Delete removes a value from the set.
//
// Parameters:
// - value: The value to be removed from the set.
func (set *Set[T]) Delete(value T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	delete(set.data, value)
}

// Clear removes all values from the set.
func (set *Set[T]) Clear() {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	// Create a new empty set with the same length as the existing set.
	set.data = make(map[T]struct{}, set.Len())
}

// OrderedSet represents an ordered set implemented using both a map and a slice.
type OrderedSet[T comparable] struct {
	mutex    sync.RWMutex
	setMap   map[T]struct{}
	setSlice []T
}

// NewOrderedSet creates a new ordered set with map and slice.
//
// Parameters:
// - capacity: The initial capacity of the ordered set.
//
// Returns:
// - A pointer to the newly created OrderedSet.
func NewOrderedSet[T comparable](capacity int) *OrderedSet[T] {
	return &OrderedSet[T]{
		setMap:   make(map[T]struct{}, capacity),
		setSlice: make([]T, 0, capacity),
	}
}

// Add adds a value to the ordered set.
//
// Parameters:
// - value: The value to be added to the ordered set.
func (orderedSet *OrderedSet[T]) Add(value T) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()
	if orderedSet.IsUnique(value) {
		orderedSet.setMap[value] = struct{}{}
		orderedSet.setSlice = append(orderedSet.setSlice, value)
	}
}

// AddAll adds multiple values to the ordered set.
//
// Parameters:
// - values: The values to be added to the ordered set.
func (orderedSet *OrderedSet[T]) AddAll(values ...T) {
	for _, value := range values {
		orderedSet.Add(value)
	}
}

// Contains checks if a value exists in the ordered set.
//
// Parameters:
// - value: The value to check for existence in the ordered set.
//
// Returns:
// - A boolean indicating whether the value is present in the ordered set.
func (orderedSet *OrderedSet[T]) Contains(value T) bool {
	orderedSet.mutex.RLock()
	defer orderedSet.mutex.RUnlock()
	_, exists := orderedSet.setMap[value]
	return exists
}

// IsUnique checks if a value is unique in the ordered set.
//
// Parameters:
// - value: The value to check for uniqueness in the ordered set.
//
// Returns:
// - A boolean indicating whether the value is unique in the ordered set.
func (orderedSet *OrderedSet[T]) IsUnique(value T) bool {
	orderedSet.mutex.RLock()
	defer orderedSet.mutex.RUnlock()
	_, exists := orderedSet.setMap[value]
	return !exists
}

// Values returns the values from the ordered set.
//
// Returns:
// - A slice containing all the values in the ordered set.
func (orderedSet *OrderedSet[T]) Values() []T {
	orderedSet.mutex.RLock()
	defer orderedSet.mutex.RUnlock()
	values := make([]T, len(orderedSet.setSlice))
	copy(values, orderedSet.setSlice)
	return values
}

// Delete removes a value from the ordered set based on the index.
//
// Parameters:
// - index: The index of the value to be removed.
func (orderedSet *OrderedSet[T]) Delete(index uint) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()
	// Check if the index is within the bounds of the slice.
	if int(index) < len(orderedSet.setSlice) {
		// Retrieve the value at the specified index.
		value := orderedSet.setSlice[index]

		// Delete from the map.
		delete(orderedSet.setMap, value)

		// Delete from the slice.
		orderedSet.setSlice = append(orderedSet.setSlice[:index], orderedSet.setSlice[index+1:]...)
	}
}

// Clear removes all values from the ordered set.
func (orderedSet *OrderedSet[T]) Clear() {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()

	// Create a new empty map with the same length as the existing map.
	orderedSet.setMap = make(map[T]struct{}, len(orderedSet.setMap))

	// Reset the slice to an empty slice, reusing the existing underlying array.
	orderedSet.setSlice = orderedSet.setSlice[:0]
}
