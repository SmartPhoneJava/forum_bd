package rerrors

import "errors"

// ErrorInvalidPath invalid path
func ErrorInvalidPath() error {
	return errors.New("Path invalid")
}

// ErrorInvalidID invalid id
func ErrorInvalidID() error {
	return errors.New("id invalid")
}

// ErrorPostConflict post conflict
func ErrorPostConflict() error {
	return errors.New("Parent post was created in another thread")
}
