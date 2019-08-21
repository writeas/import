package wfimport

import (
	"archive/zip"
	"bufio"

	"github.com/writeas/go-writeas/v2"
)

// TopLevelZipFunc return a poiter to a writeas.PostParams for any parseable
// zip.File that is not a directory. It does not traverse children.
//
// This is an example of a ZipFunc that can be used to filter files parse from
// a zip archive.
func TopLevelZipFunc(f *zip.File) (*writeas.PostParams, error) {
	if !f.FileInfo().IsDir() {
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
		post, err := fromBytes(b)
		if err != nil {
			return nil, err
		}
		return post, nil
	}
	return nil, nil
}
