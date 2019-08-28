// Copyright © 2018-2019 A Bunch Tell LLC.
//
// This is free software: you can redistribute it and/or modify
// it under the terms of the Mozilla Public License, included
// in the LICENSE file in this source code package.

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
