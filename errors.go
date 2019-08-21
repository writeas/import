package wfimport

import "errors"

var (
	// ErrEmptyFile is returned when the file is empty
	ErrEmptyFile = errors.New("file is empty")
	// ErrEmptyDir is returned when the directory is empty
	ErrEmptyDir = errors.New("directory is empty")
	// ErrInvalidFiletype is returned when the file was not plaintext
	ErrInvalidFiletype = errors.New("file type is invalid")
)
