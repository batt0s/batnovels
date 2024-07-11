package authentication

import "errors"

var (
	ErrorIncorrectPassword = errors.New("given password for given username is incorrect")
)
