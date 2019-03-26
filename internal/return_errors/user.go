package rerrors

import "errors"

// ErrorUserNotExist user not exist
func ErrorUserNotExist() error {
	return errors.New("User not exist")
}

// ErrorEmailIstaken Email is taken
func ErrorEmailIstaken() error {
	return errors.New("Email is taken")
}
