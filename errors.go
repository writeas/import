package wfimport

import "errors"

var (
	// ErrEmptyFile is returned when the file is empty
	ErrEmptyFile = errors.New("file is empty")
	// ErrInvalidContentType is returned when the detected content type does
	// not match text/* as per net/http.DetectContentType
	ErrInvalidContentType = errors.New("invalid content type")
	// ErrEmptyDir is returned when the directory is empty
	ErrEmptyDir = errors.New("directory is empty")
)
