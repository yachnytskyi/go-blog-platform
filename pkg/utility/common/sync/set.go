package sync

import "sync"

// Set represents a synchronized set implemented using a map.
type Set[T comparable] struct {
	mutex sync.Mutex
	data  map[T]struct{}
}

// NewSet creates a new set with an optional capacity.
func NewSet[T comparable](capacity uint) *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}, capacity),
	}
}

// Add adds a value to the synchronized set.
func (set *Set[T]) Add(value T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	set.data[value] = struct{}{}
}

// AddAll adds multiple values to the synchronized set.
func (set *Set[T]) AddAll(values ...T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	for _, value := range values {
		set.data[value] = struct{}{}
	}
}

// Contains checks if a value exists in the set.
// It returns true if the value is present, and false otherwise.
func (set *Set[T]) Contains(value T) bool {
	_, exists := set.data[value]
	return exists
}

// IsUnique checks if a value exists in the set.
func (set *Set[T]) IsUnique(value T) bool {
	_, exists := set.data[value]
	return !exists
}

// Len returns the number of values in the set.
// It returns true if the value does not exist in the set, and false otherwise.
func (set *Set[T]) Len() int {
	return len(set.data)
}

// Values returns the values from the synchronized set.
func (set *Set[T]) Values() []T {
	values := make([]T, 0, len(set.data))
	for value := range set.data {
		values = append(values, value)
	}
	return values
}

// Delete removes a value from the set.
func (set *Set[T]) Delete(value T) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	delete(set.data, value)
}

// Clear removes all values from the set.
func (set *Set[T]) Clear() {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	set.data = make(map[T]struct{}, len(set.data))
}

// OrderedSet represents an ordered set implemented using both a map and a slice.
type OrderedSet[T comparable] struct {
	mutex    sync.Mutex
	setMap   map[T]struct{}
	setSlice []T
}

// NewOrderedSet creates a new ordered set with map and slice.
func NewOrderedSet[T comparable](capacity uint) *OrderedSet[T] {
	return &OrderedSet[T]{
		setMap:   make(map[T]struct{}, capacity),
		setSlice: make([]T, 0, capacity),
	}
}

// Add adds a value to the ordered set.
func (orderedSet *OrderedSet[T]) Add(value T) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()

	// Check for uniqueness using the map.
	if orderedSet.IsUnique(value) {
		orderedSet.setMap[value] = struct{}{}
		orderedSet.setSlice = append(orderedSet.setSlice, value)
	}
}

// AddAll adds multiple values to the ordered set.
func (orderedSet *OrderedSet[T]) AddAll(values ...T) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()
	for _, value := range values {
		// Check for uniqueness using the map.
		if orderedSet.IsUnique(value) {
			orderedSet.setMap[value] = struct{}{}
			orderedSet.setSlice = append(orderedSet.setSlice, value)
		}
	}
}

// Contains checks if a value exists in the ordered set.
// It returns true if the value is present, and false otherwise.
func (orderedSet *OrderedSet[T]) Contains(value T) bool {
	_, exists := orderedSet.setMap[value]
	return exists
}

// IsUnique checks if a value exists in the ordered set.
// It returns true if the value does not exist in the set, and false otherwise.
func (orderedSet *OrderedSet[T]) IsUnique(value T) bool {
	_, exists := orderedSet.setMap[value]
	return !exists
}

// Len returns the number of values in the ordered set.
func (orderedSet *OrderedSet[T]) Len() int {
	return len(orderedSet.setSlice)
}

// Values returns the values from the ordered set.
func (orderedSet *OrderedSet[T]) Values() []T {
	return orderedSet.setSlice
}

// Delete removes a value from the ordered set based on the index.
func (orderedSet *OrderedSet[T]) Delete(index uint) {
	orderedSet.mutex.Lock()
	defer orderedSet.mutex.Unlock()
	// Check if the index is within the bounds of the slice and non-negative.
	deleteIndex := int(index)
	length := orderedSet.Len()
	if deleteIndex < length && length > 0 {
		// Retrieve the value at the specified index.
		value := orderedSet.setSlice[deleteIndex]

		// Delete from the map.
		delete(orderedSet.setMap, value)

		// Delete from the slice.
		orderedSet.setSlice = orderedSet.setSlice[:deleteIndex+copy(orderedSet.setSlice[deleteIndex:], orderedSet.setSlice[deleteIndex+1:])]
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
