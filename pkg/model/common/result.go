package common

type Result[T any] struct {
	Data  T
	Error error
}

// func NewResult[T any](data T, err error) *Result[T] {
// 	return &Result[T]{
// 		Data:  data,
// 		Error: err,
// 	}
// }

func NewResultWithData[T any](data T) Result[T] {
	return Result[T]{
		Data: data,
	}
}

func NewResultWithError[T any](err error) Result[T] {
	return Result[T]{
		Error: err,
	}
}
