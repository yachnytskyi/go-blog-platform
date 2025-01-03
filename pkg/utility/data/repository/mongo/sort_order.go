package mongo

import (
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

const (
	sortAscend  = 1  // Numeric value representing ascending sort order.
	sortDescend = -1 // Numeric value representing descending sort order.
)

// SetSortOrder maps a user-provided sort order string into a corresponding
// numeric value used for sorting. It handles both ascending ("ascend") and
// descending ("descend") sort orders. If the input is neither of these, it
// falls back to the default sort order.
func SetSortOrder(sortOrder string) int {
	// Map the sort order string to lower case for case-insensitive comparison.
	sortOrder = strings.ToLower(sortOrder)

	// Determine the corresponding numeric sort order value based on the input string.
	switch sortOrder {
	case constants.SortAscend:
		return sortAscend
	case constants.SortDescend:
		return sortDescend
	default:
		// If the input is neither "ascend" nor "descend", return the default sort order.
		// Check the default sort order defined in constants and return the corresponding numeric value.
		if constants.DefaultSortOrder == constants.SortAscend {
			return sortAscend
		}

		return sortDescend
	}
}
