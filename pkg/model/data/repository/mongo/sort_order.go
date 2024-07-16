package mongo

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/domain"
)

const (
	sortAscend  = 1  // Numeric value representing ascending sort order.
	sortDescend = -1 // Numeric value representing descending sort order.
)

// SetSortOrder maps a user-provided sort order string into a corresponding
// numeric value used for sorting. It handles both ascending ("ascend") and
// descending ("descend") sort orders. If the input is neither of these, it
// falls back to the default sort order defined in constants.
//
// Parameters:
//   - sortOrder: The sort order string provided by the user. Expected values are
//     "ascend" for ascending order and "descend" for descending order. The comparison
//     is case-insensitive.
//
// Returns:
// - An integer representing the sort order:
//   - 1 for ascending order ("ascend")
//   - -1 for descending order ("descend")
//   - The default sort order defined in constants if the input is not recognized.
func SetSortOrder(sortOrder string) int {
	// Map the sort order string to lower case for case-insensitive comparison.
	sortOrder = domainUtility.ToLowerString(sortOrder)

	// Determine the corresponding numeric sort order value based on the input string.
	switch sortOrder {
	case constants.SortAscend:
		// Return the numeric value for ascending sort order.
		return sortAscend
	case constants.SortDescend:
		// Return the numeric value for descending sort order.
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
