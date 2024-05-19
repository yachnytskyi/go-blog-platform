package common

// Set represents a set implemented using a map.
type Set[T comparable] map[T]struct{}

// Add adds a value to the set.
// Parameters:
// - value: The value to be added to the set.
func (set Set[T]) Add(value T) {
	set[value] = struct{}{}
}

// AddAll adds multiple values to the set.
// Parameters:
// - values: The values to be added to the set.
func (set Set[T]) AddAll(values ...T) {
	for _, value := range values {
		set[value] = struct{}{}
	}
}

// Contains checks if a value exists in the set.
// Parameters:
// - value: The value to check for existence in the set.
// Returns:
// - A boolean indicating whether the value is present in the set.
func (set Set[T]) Contains(value T) bool {
	_, exists := set[value]
	return exists
}

// IsUnique checks if a value is unique in the set.
// Parameters:
// - value: The value to check for uniqueness in the set.
// Returns:
// - A boolean indicating whether the value is unique in the set.
func (set Set[T]) IsUnique(value T) bool {
	_, exists := set[value]
	return !exists
}

// Len returns the number of values in the set.
// Returns:
// - The number of values in the set.
func (set Set[T]) Len() int {
	return len(set)
}

// Values returns the values from the set.
// Returns:
// - A slice containing all the values in the set.
func (set Set[T]) Values() []T {
	values := make([]T, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	return values
}

// Delete removes a value from the set.
// Parameters:
// - value: The value to be removed from the set.
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
	setMap   map[T]struct{} // Map to ensure uniqueness of elements.
	setSlice []T            // Slice to maintain order of elements.
}

// NewOrderedSet creates a new ordered set with a specified capacity.
// Parameters:
// - capacity: The initial capacity of the ordered set.
// Returns:
// - An OrderedSet with the specified capacity.
func NewOrderedSet[T comparable](capacity uint) OrderedSet[T] {
	return OrderedSet[T]{
		setMap:   make(map[T]struct{}, capacity),
		setSlice: make([]T, 0, capacity),
	}
}

// Add adds a value to the ordered set.
// Parameters:
// - value: The value to be added to the ordered set.
func (orderedSet *OrderedSet[T]) Add(value T) {
	// Check for uniqueness using the map.
	if orderedSet.IsUnique(value) {
		orderedSet.setMap[value] = struct{}{}
		orderedSet.setSlice = append(orderedSet.setSlice, value)
	}
}

// AddAll adds multiple values to the ordered set.
// Parameters:
// - values: The values to be added to the ordered set.
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
// Parameters:
// - value: The value to check for existence in the ordered set.
// Returns:
// - A boolean indicating whether the value is present in the ordered set.
func (orderedSet OrderedSet[T]) Contains(value T) bool {
	_, exists := orderedSet.setMap[value]
	return exists
}

// IsUnique checks if a value is unique in the ordered set.
// Parameters:
// - value: The value to check for uniqueness in the ordered set.
// Returns:
// - A boolean indicating whether the value is unique in the ordered set.
func (orderedSet OrderedSet[T]) IsUnique(value T) bool {
	_, exists := orderedSet.setMap[value]
	return !exists
}

// Values returns the values from the ordered set.
// Returns:
// - A slice containing all the values in the ordered set.
func (orderedSet OrderedSet[T]) Values() []T {
	return orderedSet.setSlice
}

// Len returns the number of values in the ordered set.
// Returns:
// - The number of values in the ordered set.
func (orderedSet OrderedSet[T]) Len() int {
	return len(orderedSet.setSlice)
}

// Delete removes a value from the ordered set based on the index.
// Parameters:
// - index: The index of the value to be removed.
func (orderedSet *OrderedSet[T]) Delete(index uint) {
	// Check if the index is within the bounds of the slice.
	deleteIndex := int(index)
	length := orderedSet.Len()
	if deleteIndex < length && deleteIndex >= 0 {
		// Retrieve the value at the specified index.
		value := orderedSet.setSlice[deleteIndex]

		// Delete from the map.
		delete(orderedSet.setMap, value)

		// Delete from the slice.
		orderedSet.setSlice = append(orderedSet.setSlice[:deleteIndex], orderedSet.setSlice[deleteIndex+1:]...)
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
