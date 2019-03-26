package rerrors

import "errors"

// ErrorVoteNotExist vote not exist
func ErrorVoteNotExist() error {
	return errors.New("Thread not exist")
}
