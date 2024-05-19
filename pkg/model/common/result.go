package common

// Result is a generic type that holds either a success value (Data) or an error (Error).
// The generic type T allows it to be used with any data type.
type Result[T any] struct {
	Data  T     // The success value of the result.
	Error error // The error value of the result, if any.
}

// NewResultOnSuccess creates a new Result instance representing a successful operation.
// It takes a value of type T and returns a Result containing that value with no error.
func NewResultOnSuccess[T any](data T) Result[T] {
	return Result[T]{
		Data: data,
	}
}

// NewResultOnFailure creates a new Result instance representing a failed operation.
// It takes an error and returns a Result containing that error with no data.
func NewResultOnFailure[T any](err error) Result[T] {
	return Result[T]{
		Error: err,
	}
}
