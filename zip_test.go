package wfimport

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type fileList []struct {
	Name, Contents string
}

var files = fileList{
	{"post.txt", "This is a post from somewhere."},
	{"books.md", "# Title of post\n\ntext body of post."},
	{"secret.txt", "shhh"},
}
var filesWDirs = fileList{
	{"post.txt", "This is a post from somewhere."},
	{"books.md", "# Title of post\n\ntext body of post."},
	{"secret.txt", "shhh"},
	{"blog/post1.txt", "some file stuff"},
	{"blog/post2.md", "shorter"},
	{"notes/test.txt", "test all the things"},
}

func TestFromZip(t *testing.T) {
	a := getTestZip(t, files)
	posts, err := FromZip(a)
	if err != nil {
		t.Fatalf("failed to get posts from archive: %v", err)
	}
	if posts == nil {
		t.Fatal("Posts was nil, expecting posts returned")
	}
	if len(posts) != len(files) {
		t.Fatalf("Post count mismatch: got %d but expected %d", len(posts), len(files))
	}
	// TODO: add check for contents, needs to update test file data above for
	// easier comparison
}

func TestTextFileZipFunc(t *testing.T) {
	a := getTestZip(t, files)
	posts, err := FromZipByFunc(a, TextFileZipFunc)
	if err != nil {
		t.Fatalf("failed to get posts from archive: %v", err)
	}
	if posts == nil {
		t.Fatal("Posts was nil, expecting posts returned")
	}
	if len(posts) != 2 { // should only be the number of files in top level, .txt
		t.Fatalf("Post count mismatch: got %d but expected %d", len(posts), 2)
	}
}

func TestFromZipDirs(t *testing.T) {
	a := getTestZip(t, filesWDirs)
	postMap, err := FromZipDirs(a)
	if err != nil {
		t.Fatalf("getting posts from zip: %v", err)
	}
	if postMap == nil {
		t.Fatal("map was nil")
	}

	if postMap["drafts"] == nil {
		t.Fatal("drafts slice should not be nil")
	}
	if len(postMap["drafts"]) != 3 {
		t.Fatalf("draft count mismatch: got %d, expecting 3", len(postMap["drafts"]))
	}
	if postMap["blog"] == nil {
		t.Fatal("blog slice should not be nil")
	}
	if len(postMap["blog"]) != 2 {
		t.Fatalf("blog count mismatch: got %d, expecting 1", len(postMap["blog"]))
	}
	if postMap["notes"] == nil {
		t.Fatal("notes slice should not be nil")
	}
	if len(postMap["notes"]) != 1 {
		t.Fatalf("notes count mismatch: got %d, expecting 1", len(postMap["notes"]))
	}
}

func getTestZip(t *testing.T, files fileList) string {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			t.Fatalf("creating file in zip: %v", err)
		}
		_, err = f.Write([]byte(file.Contents))
		if err != nil {
			t.Fatalf("writing file contents: %v", err)
		}
	}
	// add a directory, need to create file header
	err := w.Close()
	if err != nil {
		t.Fatalf("closing zip writer: %v", err)
	}

	dir := os.TempDir()
	file, err := os.Create(filepath.Join(dir, "testZip.zip"))
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	_, err = buf.WriteTo(file)
	if err != nil {
		t.Fatalf("writing to temp file: %v", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("closing temp file: %v", err)
	}

	return file.Name()
}
