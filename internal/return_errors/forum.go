package rerrors

import "errors"

// ErrorForumSlugInvalid slug invalid
func ErrorForumSlugInvalid() error {
	return errors.New("Slug invalid")
}

// ErrorForumSlugIsTaken slug is taken
func ErrorForumSlugIsTaken() error {
	return errors.New("Slug is taken")
}

// ErrorForumNotExist forum not exist
func ErrorForumNotExist() error {
	return errors.New("Forum not exist")
}
