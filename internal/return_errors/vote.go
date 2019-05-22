package rerrors

import "errors"

// ErrorVoteNotExist vote not exist
func ErrorVoteNotExist() error {
	return errors.New("Thread not exist")
}

// ErrorVoteNotExist vote not exist
func ErrorVoteInvalidAuthor() error {
	return errors.New("No such user")
}
