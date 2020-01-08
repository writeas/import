// Copyright © 2019-2020 A Bunch Tell LLC. and contributors.
//
// This is free software: you can redistribute it and/or modify
// it under the terms of the Mozilla Public License, included
// in the LICENSE file in this source code package.

package wfimport

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFromDirectory(t *testing.T) {
	testDir := "test"
	files := []string{"test.md", "test2.txt", "test3"}
	dirs := []string{"one", "two.md"}

	// set up test directory and children
	err := os.Mkdir(testDir, os.ModeDir|os.ModePerm)
	defer os.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("failed to create base test dir: %v", err)
	}
	for _, fn := range files {
		f, err := os.Create(filepath.Join(testDir, fn))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		if strings.HasSuffix(fn, ".md") {
			_, err = f.WriteString(fmt.Sprintf("# a title\n%s", fn))
		} else {
			_, err = f.WriteString(fn)
		}
		if err != nil {
			t.Fatalf("failed to write test file contents: %v", err)
		}
		f.Close()
	}
	for _, dn := range dirs {
		err := os.Mkdir(filepath.Join(testDir, dn), os.ModeDir)
		if err != nil {
			t.Fatalf("failed to create test dir: %v", err)
		}
	}

	posts, err := FromDirectory(testDir)
	if err != nil {
		t.Fatalf("failed to parse files from directory: %v", err)
	}

	numExpected := len(files)
	numResults := len(posts)
	if numResults != numExpected {
		t.Fatalf("post count mismatch.\bgot: %d\nexpecting: %d", numResults, numExpected)
	}
}

func TestFromDirectoryMatch(t *testing.T) {
	testDir := "test"
	files := []string{"test.md", "test2.txt", "test3"}

	// set up test directory and children
	err := os.Mkdir(testDir, os.ModeDir|os.ModePerm)
	defer os.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("failed to create base test dir: %v", err)
	}
	for _, fn := range files {
		f, err := os.Create(filepath.Join(testDir, fn))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		_, err = f.WriteString(fn)
		if err != nil {
			t.Fatalf("failed to write test file contents: %v", err)
		}
		f.Close()
	}

	posts, err := FromDirectoryMatch(testDir, `\d+`)
	if err != nil {
		t.Fatalf("failed to parse files from directory: %v", err)
	}

	numExpected := 2 // length of files matching expression
	numResults := len(posts)
	if numResults != numExpected {
		t.Fatalf("post count mismatch.\bgot: %d\nexpecting: %d", numResults, numExpected)
	}

	posts, err = FromDirectoryMatch(testDir, `test`)
	if err != nil {
		t.Fatalf("failed to parse files from directory: %v", err)
	}

	numExpected = 3 // length of files matching expression
	numResults = len(posts)
	if numResults != numExpected {
		t.Fatalf("post count mismatch.\bgot: %d\nexpecting: %d", numResults, numExpected)
	}
}

func TestFromDirectoryNoPath(t *testing.T) {
	posts, err := FromDirectory("")
	if err == nil {
		t.Fatal("error is nil but should not open directory without name")
	}
	if posts != nil {
		t.Fatal("posts returned but should be nil")
	}
}

func TestFromDirectoryRelativePath(t *testing.T) {
	testDir := "test"
	err := os.Mkdir(testDir, os.ModeDir|os.ModePerm)
	if err != nil {
		t.Fatal("failed to create test dir")
	}
	defer os.Remove(testDir)

	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("failed to change into test dir: %v", err)
	}
	defer os.Chdir("../")

	f, err := os.Create("blog.md")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() {
		f.Close()
		os.Remove("blog.md")
	}()
	_, err = f.WriteString("hello")
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}

	post, err := FromDirectory(".")
	if err != nil {
		t.Fatalf("failed to parse from relative directory: %v", err)
	}

	if post == nil {
		t.Fatal("post is nil but should exist")
	}
}

func TestFromDirectoryEmptyDir(t *testing.T) {
	testDir := "test"
	err := os.Mkdir(testDir, os.ModeDir|os.ModePerm)
	if err != nil {
		t.Fatal("failed to create test dir")
	}
	defer os.Remove(testDir)

	posts, err := FromDirectory(testDir)
	if err == nil {
		t.Fatal("error is nil but directory was empty")
	}
	if err != ErrEmptyDir {
		t.Fatalf("wrong error returned:\ngot %v\nexpecting%v", err, ErrEmptyDir)
	}
	if posts != nil {
		t.Fatal("posts returned but should be nil")
	}
}

func TestFromDirectoryErrors(t *testing.T) {
	testDir := "test"
	files := []string{"test.md", "test2.txt", "test3"}

	// set up test directory and children
	err := os.Mkdir(testDir, os.ModeDir|os.ModePerm)
	defer os.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("failed to create base test dir: %v", err)
	}
	for _, fn := range files {
		f, err := os.Create(filepath.Join(testDir, fn))
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		f.Close()
	}
	posts, err := FromDirectory(testDir)
	if err == nil {
		t.Fatal("error was nil but no files have contents")
	}
	if len(posts) == len(files) {
		t.Fatal("files with errors were returned, should be skipped")
	}
}

func TestFromFile(t *testing.T) {
	filename := "test.txt"
	postBody := `test post
	
	this is a test`

	// create file to read
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer os.Remove(filename)

	_, err = file.WriteString(postBody)
	if err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	// read from file
	post, err := FromFile(filename)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	if post == nil {
		t.Fatal("post was nil")
	}
	title, body := extractTitle(postBody)
	if post.Content != body {
		t.Logf("post content mismatch.")
		t.Logf("got:\n%s", post.Content)
		t.Logf("expected:\n%s", body)
		t.FailNow()
	}
	if post.Title != title {
		t.Logf("post title mismatch.")
		t.Logf("got:\n%s", post.Title)
		t.Logf("expected:\n%s", title)
		t.FailNow()
	}
}

func TestFromFileNoPath(t *testing.T) {
	post, err := FromFile("")
	if err == nil {
		t.Fatal("error was nil but should not open an empty path")
	}
	if post != nil {
		t.Fatal("post was returned but should be nil")
	}
}

func TestFromFileEmptyBody(t *testing.T) {
	testFile := "test.txt"
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	post, err := FromFile(testFile)
	if err != ErrEmptyFile {
		t.Fatalf("error was nil but should error with empty file")
	}
	if post != nil {
		t.Fatal("post was returned but should be nil")
	}
}

func TestFromBytes(t *testing.T) {
	tt := []struct {
		Name  string
		Bytes []byte
		Error error
	}{
		{
			Name:  "text file",
			Bytes: []byte("This is just some text."),
			Error: nil,
		}, {
			Name:  "markdown file",
			Bytes: []byte("# header\n\nText and more text."),
			Error: nil,
		}, {
			Name:  "jpeg image",
			Bytes: []byte("\xFF\xD8\xFF"),
			Error: ErrInvalidContentType,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			p, err := fromBytes(tc.Bytes)
			if err != tc.Error {
				t.Fatalf("got error %v but expected error %v", err, tc.Error)
			}

			if tc.Error == nil && p == nil {
				t.Fatal("got back nil post and nil error")
			}

			if tc.Error == nil && p.Content == "" {
				t.Fatal("got back empty content")
			}
		})
	}
}
