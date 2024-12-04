package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

func TestIsTimeNotValidPastTime(t *testing.T) {
	t.Parallel()
	PastTime := time.Now().Add(-100 * time.Millisecond)
	result := validator.IsTimeNotValid(PastTime)

	assert.True(t, result, test.NotFailureMessage)
}

func TestIsTimeNotValidCurrentTime(t *testing.T) {
	t.Parallel()
	currentTime := time.Now()
	result := validator.IsTimeNotValid(currentTime)

	assert.True(t, result, test.NotFailureMessage)
}

func TestIsTimeNotValidFutureTime(t *testing.T) {
	t.Parallel()
	futureTime := time.Now().Add(100 * time.Millisecond)
	result := validator.IsTimeNotValid(futureTime)

	assert.False(t, result, test.FailureMessage)
}
