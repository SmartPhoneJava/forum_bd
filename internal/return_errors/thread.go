package rerrors

import "errors"

// ErrorThreadNotExist Thread not exist
func ErrorThreadNotExist() error {
	return errors.New("Thread not exist")
}

// ErrorThreadConflict Thread is taken
func ErrorThreadConflict() error {
	return errors.New("Thread is taken")
}
