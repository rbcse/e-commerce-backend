package error

import "errors"

var (
	CustomerEmailAlreadyExists = errors.New("customer email already exists")
	CustomerPhoneNumberAlreadyExists = errors.New("customer phone number already exists")
)