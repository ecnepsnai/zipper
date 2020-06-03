# zipper

[![Go Report Card](https://goreportcard.com/badge/github.com/ecnepsnai/zipper?style=flat-square)](https://goreportcard.com/report/github.com/ecnepsnai/zipper)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/ecnepsnai/zipper)
[![Releases](https://img.shields.io/github/release/ecnepsnai/zipper/all.svg?style=flat-square)](https://github.com/ecnepsnai/zipper/releases)
[![LICENSE](https://img.shields.io/github/license/ecnepsnai/zipper.svg?style=flat-square)](https://github.com/ecnepsnai/zipper/blob/master/LICENSE)

zipper provides a simple interface to create zip archives in a Go application.

Yup, that's it. It's a one-trick-pony.

# Installation

```
go get github.com/ecnepsnai/zipper
```

# Usage

## Create a new ZIP Archive

```go
import (
	"github.com/ecnepsnai/zipper"
)

func Example() {
	archive, err := zipper.New("/path/to/your/zip/file")
	if err != nil {
		panic(err)
	}

	// Add your files
	// (See examples below)

	archive.Close()
}
```

## Adding Files

There are three ways to add files to the archive: Using a writer, copying bytes, and inserting an existing file.

### Using a Writer

`Archive.NewFile` takes in a file name and returns a `io.Writer`. Writing to that writer will add data to the file name
specified in the zip archive.

```go
import (
	"bytes"
	"io"

	"github.com/ecnepsnai/zipper"
)

func Example() {
	var archive *zip.Archive // See zip.NewZipFile for instructions

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
```

### Copying bytes

`Archive.AddFile` takes in a file name and slice of bytes. A new file is created with the given file name and populated
with the bytes slice in the zip archive.

```go
import (
	"github.com/ecnepsnai/zipper"
)

func Example() {
	var archive *zip.Archive // See zip.NewZipFile for instructions

	data := []byte("Hello world!")

	// AddFile adds a new file to the zip archive with the provided file name populated with the given data
	if err := archive.AddFile("my_file.txt", data); err != nil {
		panic(err)
	}
}
```

### Inserting an Existing File

`Archive.InsertFile` takes in a file path and will attempt to copy that file to a new file in the zip archive. It will
only use the file name, stripping any parent directories.

```go
import (
	"github.com/ecnepsnai/zipper"
)

func Example() {
	var archive *zip.Archive // See zip.NewZipFile for instructions

	path := "/path/to/existing/file"

	// InsertFile copies an existing local file to the zip archive
	if err := archive.InsertFile(path); err != nil {
		panic(err)
	}
}
```