package zipper_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"path"
	"testing"

	"github.com/ecnepsnai/zipper"
)

var tmpDir string

func testSetup() {
	tmp, err := ioutil.TempDir("", "zip")
	if err != nil {
		panic(err)
	}
	tmpDir = tmp
}

func testTeardown() {
	os.RemoveAll(tmpDir)
}

func randomFilePath(ext string) string {
	n := mrand.Intn(99999-10000) + 10000
	return path.Join(tmpDir, fmt.Sprintf("zip_%d.%s", n, ext))
}

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	testTeardown()
	os.Exit(retCode)
}

func TestZipFile(t *testing.T) {
	zipPath := randomFilePath("zip")
	archive, err := zipper.New(zipPath)
	if err != nil {
		t.Fatalf("Error making zip archive '%s': %s", zipPath, err.Error())
	}

	// Add file (writer)
	data := []byte("Hello world!")
	writer, err := archive.NewFile("writer.txt")
	if err != nil {
		t.Fatalf("Error opening writer for new file in archive: %s", err.Error())
	}
	l, err := io.Copy(writer, bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Error copying data to archive file writer: %s", err.Error())
	}
	if int64(len(data)) != l {
		t.Fatalf("Written data does not match actual length (%d != %d)", len(data), l)
	}

	// Add file (data)
	if err := archive.AddFile("data.txt", data); err != nil {
		t.Fatalf("Error adding data to archive: %s", err.Error())
	}

	// Insert existing file
	filePath := randomFilePath("txt")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Error opening temporary file for writing: %s", err.Error())
	}
	if _, err := io.Copy(f, bytes.NewBuffer(data)); err != nil {
		t.Fatalf("Error copying data to temporary file: %s", err.Error())
	}
	f.Close()
	if err := archive.InsertFile(filePath); err != nil {
		t.Fatalf("Error inserting existing file into archive: %s", err.Error())
	}

	archive.Close()
}

func ExampleNew() {
	archive, err := zipper.New("/path/to/your/zip/file")
	if err != nil {
		panic(err)
	}

	// Add your files
	// (See examples below)

	archive.Close()
}

func ExampleArchive_NewFile() {
	var archive *zipper.Archive // See zipper.NewZipFile for instructions

	// NewFile adds a new file to the zip archive and returns a writer
	writer, err := archive.NewFile("my_file.txt")
	if err != nil {
		panic(err)
	}

	// NewFile returns a writer, so you can use any of the writer operations on it
	// such as io.Copy
	if _, err := io.Copy(writer, bytes.NewBuffer([]byte("Hello world!"))); err != nil {
		panic(err)
	}
}

func ExampleArchive_AddFile() {
	var archive *zipper.Archive // See zipper.NewZipFile for instructions

	data := []byte("Hello world!")

	// AddFile adds a new file to the zip archive with the provided file name populated with the given data
	if err := archive.AddFile("my_file.txt", data); err != nil {
		panic(err)
	}
}

func ExampleArchive_InsertFile() {
	var archive *zipper.Archive // See zipper.NewZipFile for instructions

	path := "/path/to/existing/file"

	// InsertFile copies an existing local file to the zip archive
	if err := archive.InsertFile(path); err != nil {
		panic(err)
	}
}
