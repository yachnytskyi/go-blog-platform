package common

// Set represents a set implemented using a map.
type Set[T comparable] map[T]struct{}

func (set Set[T]) Add(value T) {
	set[value] = struct{}{}
}

func (set Set[T]) AddAll(values ...T) {
	for _, value := range values {
		set.Add(value)
	}
}

func (set Set[T]) Contains(value T) bool {
	_, exists := set[value]
	return exists
}

func (set Set[T]) IsUnique(value T) bool {
	_, exists := set[value]
	if exists {
		return false
	}

	return true
}

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) Values() []T {
	values := make([]T, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	return values
}

func (set Set[T]) Delete(value T) {
	delete(set, value)
}

func (set *Set[T]) Clear() {
	*set = make(Set[T], len(*set))
}

// OrderedSet represents an ordered set implemented using both a map and a slice.
type OrderedSet[T comparable] struct {
	setMap   map[T]struct{} // Map to ensure uniqueness of elements.
	setSlice []T            // Slice to maintain order of elements.
}

func NewOrderedSet[T comparable](capacity uint) OrderedSet[T] {
	return OrderedSet[T]{
		setMap:   make(map[T]struct{}, capacity),
		setSlice: make([]T, 0, capacity),
	}
}

func (orderedSet *OrderedSet[T]) Add(value T) {
	// Check for uniqueness using the map.
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

func (orderedSet OrderedSet[T]) Contains(value T) bool {
	_, exists := orderedSet.setMap[value]
	return exists
}

func (orderedSet OrderedSet[T]) IsUnique(value T) bool {
	_, exists := orderedSet.setMap[value]
	if exists {
		return false
	}

	return true
}

func (orderedSet OrderedSet[T]) Values() []T {
	return orderedSet.setSlice
}

func (orderedSet *OrderedSet[T]) Delete(index uint) {
	if int(index) < len(orderedSet.setSlice) {
		value := orderedSet.setSlice[index]
		delete(orderedSet.setMap, value)
		orderedSet.setSlice = append(orderedSet.setSlice[:index], orderedSet.setSlice[index+1:]...)
	}
}

func (orderedSet *OrderedSet[T]) Clear() {
	orderedSet.setMap = make(map[T]struct{}, len(orderedSet.setMap))
	orderedSet.setSlice = orderedSet.setSlice[:0]
}
