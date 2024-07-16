package common

// Result is a generic type that holds either a success value (Data) or an error (Error).
// The generic type T allows it to be used with any data type.
// It provides a standard way to represent the outcome of an operation,
// making it easier to handle success and failure cases consistently.
//
// Fields:
// - Data: The success value of the result. This field is populated only in case of success.
// - Error: The error value of the result. This field is populated only in case of failure.
type Result[T any] struct {
	Data  T     // The success value of the result.
	Error error // The error value of the result, if any.
}

// NewResultOnSuccess creates a new Result instance representing a successful operation.
// It initializes the Data field with the provided value and leaves the Error field nil.
// This function is used when an operation completes successfully.
//
// Parameters:
// - data: The success value of type T to be contained in the Result.
//
// Returns:
// - Result[T]: The result containing the success value and no error.
func NewResultOnSuccess[T any](data T) Result[T] {
	return Result[T]{
		Data: data,
	}
}

// NewResultOnFailure creates a new Result instance representing a failed operation.
// It initializes the Error field with the provided error and leaves the Data field with its zero value.
// This function is used when an operation encounters an error.
//
// Parameters:
// - err: The error to be contained in the Result.
//
// Returns:
// - Result[T]: The result containing the error and no success value.
func NewResultOnFailure[T any](err error) Result[T] {
	return Result[T]{
		Error: err,
	}
}
