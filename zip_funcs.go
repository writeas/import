// Copyright © 2019-2020 A Bunch Tell LLC. and contributors.
//
// This is free software: you can redistribute it and/or modify
// it under the terms of the Mozilla Public License, included
// in the LICENSE file in this source code package.

package wfimport

import (
	"archive/zip"
	"bufio"
	"path/filepath"
	"strings"

	"github.com/writeas/go-writeas/v2"
)

// TopLevelZipFunc returns a pointer to a writeas.PostParams for any parseable
// zip.File that is not a directory. It does not traverse children.
//
// This is an example of a ZipFunc that can be used to filter files parse from
// a zip archive.
func TopLevelZipFunc(f *zip.File) (*writeas.PostParams, error) {
	if !f.FileInfo().IsDir() {
		return openAndParse(f)
	}
	return nil, nil
}

// TextFileZipFunc parses .txt files into PostParams
func TextFileZipFunc(f *zip.File) (*writeas.PostParams, error) {
	if !f.FileInfo().IsDir() && filepath.Ext(f.FileHeader.Name) == ".txt" {
		return openAndParse(f)
	}
	return nil, nil
}

func openAndParse(f *zip.File) (*writeas.PostParams, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	r := bufio.NewReader(rc)
	b := make([]byte, f.FileInfo().Size())
	_, err = r.Read(b)
	if err != nil {
		return nil, err
	}
	p, err := fromBytes(b)
	if err != nil {
		return nil, err
	}
	p.Created = &f.Modified

	p.ID, p.Slug, p.Collection = filenameParts(f.FileHeader.Name)
	return p, nil
}

func filenameParts(filename string) (id, slug, coll string) {
	filename = strings.TrimSuffix(filename, ".txt")
	seg := strings.Split(filename, "/")
	if len(seg) > 1 {
		coll = seg[0]
		filename = seg[1]
	}
	seg = strings.Split(filename, "_")
	if len(seg) > 1 {
		slug = seg[0]
		filename = seg[1]
	}
	id = filename
	return
}
