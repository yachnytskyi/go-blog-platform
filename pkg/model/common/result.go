package common

// Result is a generic type that holds either a success value (Data) or an error (Error).
type Result[T any] struct {
	Data  T     // The success value of the result.
	Error error // The error value of the result, if any.
}

func NewResultOnSuccess[T any](data T) Result[T] {
	return Result[T]{
		Data: data,
	}
}

func NewResultOnFailure[T any](err error) Result[T] {
	return Result[T]{
		Error: err,
	}
}

func (r Result[T]) IsError() bool {
	return r.Error != nil
}
