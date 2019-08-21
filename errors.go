package wfimport

import "errors"

var (
	// ErrEmptyFile is returned when the file is empty
	ErrEmptyFile = errors.New("file is empty")
	// ErrEmptyDir is returned when the directory is empty
	ErrEmptyDir = errors.New("directory is empty")
)
