package wfimport

import (
	"archive/zip"
	"path/filepath"

	"github.com/writeas/go-writeas"
)

// ZipCollections holds a map of collections of post params.
// The keys are the collection name, parsed from the directory structure.
// Draft posts are included under the key drafts, those that were top level
// files in the archive.
type ZipCollections map[string][]*writeas.PostParams

// ZipFunc should return a pointer to a writeas.PostParams for any zip.File
// that meets criteria. It is used in FromZipFunc to filter the archives files.
//
// For an example, see the file zip_funcs.go
type ZipFunc func(f *zip.File) (*writeas.PostParams, error)

// FromZip opens a zip archive and returns a slice of *writeas.PostParams
// and an error if any. It only reads the top level of the archive tree.
func FromZip(archive string) ([]*writeas.PostParams, error) {
	return FromZipByFunc(archive, topLevelZipFunc)
}

// FromZipByFunc opens an archive and filters the contents according to the
// passed ZipFunc. It returns a slice of writeas.PostParams and any error.
func FromZipByFunc(archive string, f ZipFunc) ([]*writeas.PostParams, error) {
	a, err := zip.OpenReader(archive)
	if err != nil {
		return nil, err
	}
	defer a.Close()

	return postsFromZipFiles(a.File, f)
}

// FromZipDirs opens a zip archive and returns a map of post collections
// and an error if any.
//
// The map is of [string][]*writeas.PostParams where the string key is the name
// of the directory. The top level directory posts will be 'drafts'.
func FromZipDirs(archive string) (ZipCollections, error) {
	return FromZipDirsByFunc(archive, topLevelZipFunc)
}

// FromZipDirsByFunc works as FromZipDirs but filtering files through f.
func FromZipDirsByFunc(archive string, f ZipFunc) (ZipCollections, error) {
	return postsFromZipDirs(archive, f)
}

func postsFromZipFiles(files []*zip.File, f ZipFunc) ([]*writeas.PostParams, error) {
	posts := []*writeas.PostParams{}
	for _, file := range files {
		post, err := f(file)
		if err == ErrEmptyFile {
			continue
		} else if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if len(posts) > 0 {
		return posts, nil
	}
	return nil, nil
}

func postsFromZipDirs(archive string, f ZipFunc) (ZipCollections, error) {
	out := make(ZipCollections)
	a, err := zip.OpenReader(archive)
	if err != nil {
		return nil, err
	}
	defer a.Close()

	drafts := []*writeas.PostParams{}
	dirs := make(map[string][]*zip.File)
	for _, file := range a.File {
		dir, _ := filepath.Split(file.Name)
		if dir != "" {
			dir = filepath.Dir(dir)
			if dirs[dir] == nil {
				dirs[dir] = []*zip.File{}
			}
			dirs[dir] = append(dirs[dir], file)
		} else {
			post, err := f(file)
			if err == ErrEmptyFile {
				continue
			} else if err != nil {
				return nil, err
			}
			drafts = append(drafts, post)
		}
	}
	out["drafts"] = drafts
	for dirName, dirList := range dirs {
		if out[dirName] == nil {
			out[dirName] = []*writeas.PostParams{}
		}
		for _, file := range dirList {
			post, err := f(file)
			if err == ErrEmptyFile {
				continue
			} else if err != nil {
				return nil, err
			}
			out[dirName] = append(out[dirName], post)
		}
	}

	return out, nil
}
