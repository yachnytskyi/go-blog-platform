package common

type Result[T any] struct {
	Data  T
	Error error
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
