package validator

func IsErrorNil(err error) bool {
	return err == nil
}

func IsErrorNotNil(err error) bool {
	return err != nil
}
