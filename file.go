// Copyright © 2018-2019 A Bunch Tell LLC.
//
// This is free software: you can redistribute it and/or modify
// it under the terms of the Mozilla Public License, included
// in the LICENSE file in this source code package.

package wfimport

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/writeas/go-writeas/v2"
)

// FromDirectoryMatch reads all text and markdown files in path that match the
// pattern returning the parsed posts and an error if any.
//
// The pattern should be a valid regex, for more details see
// https://golang.org/s/re2syntax or run `go doc regexp/syntax`
func FromDirectoryMatch(path, pattern string) ([]*writeas.PostParams, error) {
	return fromDirectory(path, pattern)
}

// FromDirectory reads all text and markdown files in path and returns the
// parsed posts and an error if any.
func FromDirectory(path string) ([]*writeas.PostParams, error) {
	return fromDirectory(path, "")
}

// fromDirectory takes an 'optional' pattern, if an empty string is passed
// all valid txt and md files will be included under path.
// Otherwise pattern should be a valid regex per MatchFromDirectory
func fromDirectory(path, pattern string) ([]*writeas.PostParams, error) {
	if pattern == "" {
		pattern = "."
	}
	rx, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrEmptyDir
	}

	var postErrors error
	posts := []*writeas.PostParams{}
	for _, f := range list {
		if !f.IsDir() {
			filename := f.Name()
			if rx.MatchString(filename) {
				post, err := FromFile(filepath.Join(path, filename))
				if err != nil {
					postErrors = multierror.Append(postErrors, err)
					continue
				}

				posts = append(posts, post)
			}
		}
	}
	return posts, postErrors
}

// FromFile reads in a file from path and returns the parsed post and an error
// if any. The title will be extracted from the first markdown level 1 header.
func FromFile(path string) (*writeas.PostParams, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	p, err := fromBytes(b)
	if err != nil {
		return nil, err
	}
	created := info.ModTime()
	p.Created = &created

	return p, nil
}

func fromBytes(b []byte) (*writeas.PostParams, error) {
	if len(b) == 0 {
		return nil, ErrEmptyFile
	}

	if contentType := http.DetectContentType(b); !strings.HasPrefix(contentType, "text/") {
		return nil, ErrInvalidContentType
	}

	title, body := extractTitle(string(b))
	post := writeas.PostParams{
		Title:   title,
		Content: body,
	}

	return &post, nil
}

// TODO: copied from writeas/web-core/posts due to errors with package imports
// maybe also find a way to grab the first line as a title in plain text files
func extractTitle(content string) (title string, body string) {
	if hashIndex := strings.Index(content, "# "); hashIndex == 0 {
		eol := strings.IndexRune(content, '\n')
		// First line should start with # and end with \n
		if eol != -1 {
			body = strings.TrimLeft(content[eol:], " \t\n\r")
			title = content[len("# "):eol]
			return
		}
	}
	body = content
	return
}
