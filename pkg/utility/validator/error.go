package validator

func IsError(err error) bool {
	return err != nil
}
