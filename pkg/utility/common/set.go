package common

// Set represents a set implemented using a map.
type Set[T comparable] map[T]struct{}

// Add adds a value to the set.
func (set Set[T]) Add(value T) {
	set[value] = struct{}{}
}

// AddAll adds multiple values to the set.
func (set Set[T]) AddAll(values ...T) {
	for _, value := range values {
		set[value] = struct{}{}
	}
}

// Contains checks if a value exists in the set.
// It returns true if the value is present, and false otherwise.
func (set Set[T]) Contains(value T) bool {
	_, exists := set[value]
	return exists
}

// IsUnique checks if a value exists in the set.
// It returns true if the value does not exist in the set, and false otherwise.
func (set Set[T]) IsUnique(value T) bool {
	_, exists := set[value]
	return !exists
}

// Len returns the number of values in the set.
func (set Set[T]) Len() int {
	return len(set)
}

// Values returns the values from the set.
func (set Set[T]) Values() []T {
	values := make([]T, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	return values
}

// Delete removes a value from the set.
func (set Set[T]) Delete(value T) {
	delete(set, value)
}

// Clear removes all values from the set.
func (set *Set[T]) Clear() {
	if set.Len() > 0 {
		// Create a new empty map with the same length as the existing map.
		*set = make(map[T]struct{}, set.Len())
	}
}

// OrderedSet represents an ordered set implemented using both a map and a slice.
type OrderedSet[T comparable] struct {
	setMap   map[T]struct{}
	setSlice []T
}

// NewOrderedSet creates a new ordered set with map and slice.
func NewOrderedSet[T comparable](capacity uint) OrderedSet[T] {
	return OrderedSet[T]{
		setMap:   make(map[T]struct{}, capacity),
		setSlice: make([]T, 0, capacity),
	}
}

// Add adds a value to the ordered set.
func (orderedSet *OrderedSet[T]) Add(value T) {
	// Check for uniqueness using the map.
	if orderedSet.IsUnique(value) {
		orderedSet.setMap[value] = struct{}{}
		orderedSet.setSlice = append(orderedSet.setSlice, value)
	}
}

// AddAll adds multiple values to the ordered set.
func (orderedSet *OrderedSet[T]) AddAll(values ...T) {
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
func (orderedSet OrderedSet[T]) Contains(value T) bool {
	_, exists := orderedSet.setMap[value]
	return exists
}

// IsUnique checks if a value exists in the ordered set.
// It returns true if the value does not exist in the set, and false otherwise.
func (orderedSet OrderedSet[T]) IsUnique(value T) bool {
	_, exists := orderedSet.setMap[value]
	return !exists
}

// Values returns the values from the ordered set.
func (orderedSet OrderedSet[T]) Values() []T {
	return orderedSet.setSlice
}

// Len returns the number of values in the ordered set.
func (orderedSet OrderedSet[T]) Len() int {
	return len(orderedSet.setSlice)
}

// Delete removes a value from the ordered set based on the index.
func (orderedSet *OrderedSet[T]) Delete(index uint) {
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
	if orderedSet.Len() > 0 {
		// Create a new empty map with the same length as the existing map.
		orderedSet.setMap = make(map[T]struct{}, len(orderedSet.setMap))

		// Reset the slice to an empty slice, reusing the existing underlying array.
		orderedSet.setSlice = orderedSet.setSlice[:0]
	}
}
