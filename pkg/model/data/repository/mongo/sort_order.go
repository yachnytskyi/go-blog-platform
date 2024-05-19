package mongo

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
)

const (
	sortAscend  = 1  // Numeric value representing ascending sort order.
	sortDescend = -1 // Numeric value representing descending sort order.
)

// SetSortOrder converts a user-provided sort order string into a corresponding
// numeric value used for sorting. It handles both ascending ("ascend") and
// descending ("descend") sort orders. If the input is neither of these, it
// falls back to the default sort order defined in constants.
// Parameters:
// - sortOrder: The sort order string provided by the user.
// Returns:
// - An integer representing the sort order, where 1 means ascending and -1 means descending.
func SetSortOrder(sortOrder string) int {
	sortOrder = domain.ToLowerString(sortOrder)
	switch sortOrder {
	case constants.SortAscend:
		return sortAscend
	case constants.SortDescend:
		return sortDescend
	default:
		if constants.DefaultSortOrder == constants.SortAscend {
			return sortAscend
		}

		return sortDescend
	}
}
