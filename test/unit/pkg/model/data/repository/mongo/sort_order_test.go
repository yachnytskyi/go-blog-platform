package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/data/repository/mongo"
)

const (
	sortAscend  = 1
	sortDescend = -1
)

func TestSetSortOrderAscend(t *testing.T) {
	t.Parallel()

	result := utility.SetSortOrder(constants.SortAscend)
	assert.Equal(t, sortAscend, result)
}

func TestSetSortOrderDescend(t *testing.T) {
	t.Parallel()

	result := utility.SetSortOrder(constants.SortDescend)
	assert.Equal(t, sortDescend, result)
}

func TestSetSortOrderInvalid(t *testing.T) {
	t.Parallel()

	result := utility.SetSortOrder("invalid")
	if constants.DefaultSortOrder == constants.SortAscend {
		assert.Equal(t, sortAscend, result)

	}

	assert.Equal(t, sortDescend, result)
}

func TestSetSortOrderEmpty(t *testing.T) {
	t.Parallel()

	result := utility.SetSortOrder("")
	if constants.DefaultSortOrder == constants.SortAscend {
		assert.Equal(t, sortAscend, result)
	}

	assert.Equal(t, sortDescend, result)
}

func TestSetSortOrderMixedCase(t *testing.T) {
	t.Parallel()

	result := utility.SetSortOrder("AscEnD")
	assert.Equal(t, sortAscend, result)
}

func TestSetSortOrderDescendMixedCase(t *testing.T) {
	t.Parallel()

	result := utility.SetSortOrder("DeScEnD")
	assert.Equal(t, sortDescend, result)
}
