// Package zipper provides a simple interface to create zip archives in a Go application.
package zipper

import (
	// For "golang zip, not GNU zip"
	gzip "archive/zip"
	"io"
	"os"
)

// Archive describes a zip archive
type Archive struct {
	FilePath  string
	zipFile   *os.File
	zipWriter *gzip.Writer
}

// New create a new empty zip archive with the given file path
func New(FilePath string) (*Archive, error) {
	newfile, err := os.Create(FilePath)
	if err != nil {
		return nil, err
	}

	zipWriter := gzip.NewWriter(newfile)

	archive := Archive{
		FilePath:  FilePath,
		zipFile:   newfile,
		zipWriter: zipWriter,
	}

	return &archive, nil
}

// NewFile takes in a file name and returns a `io.Writer`. Writing to that writer will add data to the file name
// specified in the zip archive.
func (a *Archive) NewFile(name string) (io.Writer, error) {
	return a.zipWriter.Create(name)
}

// AddFile takes in a file name and slice of bytes. A new file is created with the given file name and populated
// with the bytes slice in the zip archive.
func (a *Archive) AddFile(name string, data []byte) error {
	writer, err := a.zipWriter.Create(name)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

// InsertFile takes in a file path and will attempt to copy that file to a new file in the zip archive. It will
// only use the file name, stripping any parent directories.
func (a *Archive) InsertFile(path string) error {
	zipfile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Get the file information
	info, err := zipfile.Stat()
	if err != nil {
		return err
	}

	header, err := gzip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = gzip.Deflate

	writer, err := a.zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err := io.Copy(writer, zipfile); err != nil {
		return err
	}

	return nil
}

// Close flush and save the archive
func (a *Archive) Close() {
	a.zipWriter.Close()
	a.zipFile.Close()
}
