package mongo

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
)

const (
	sortAscend  = 1
	sortDescend = -1
)

// SetSortOrder converts a user-provided sort order string into a corresponding
// numeric value used for sorting. It handles both ascending ("ascend") and
// descending ("descend") sort orders. If the input is neither of these, it
// falls back to the default sort order defined in constants.
func SetSortOrder(sortOrder string) int {
	sortOrder = domainUtility.ToLowerString(sortOrder)
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
