package rerrors

import "errors"

// ErrorInvalidLimit invalid limit
func ErrorInvalidLimit() error {
	return errors.New("Invalid limit")
}

// ErrorInvalidDate invalid date
func ErrorInvalidDate() error {
	return errors.New("Invalid date")
}

// ErrorInvalidName Invalid name
func ErrorInvalidName() error {
	return errors.New("Invalid name")
}

// ErrorNoBody call it, if client
// didnt send you body, when you need it
func ErrorNoBody() error {
	return errors.New("Cant found parameters")
}
