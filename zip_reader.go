package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"errors"
	"strings"
)

// ZipReader is the entry point for using an archive reader on a Rar archive
var ZipReader zipFormatReader

type zipFormatEntry struct {
	zipEntry *zip.File
}

func (entry zipFormatEntry) Name() string {
	if entry.zipEntry != nil {
		return entry.zipEntry.Name
	}
	return ""
}

func (entry zipFormatEntry) IsDirectory() bool {
	// just the suffix of '/' should be enough
	return (entry.zipEntry.CompressedSize64 == 0 && entry.zipEntry.UncompressedSize64 == 0 && strings.HasSuffix(entry.Name(), "/"))
}

func (entry zipFormatEntry) Mode() os.FileMode {
	if entry.zipEntry != nil {
		return entry.zipEntry.FileInfo().Mode()
	}
	return os.ModeAppend
}

func (entry *zipFormatEntry) Write(output io.Writer) error {
	if entry.zipEntry == nil {
		return errors.New("no Reader")
	}
	rc, err := entry.zipEntry.Open()
	if err != nil {
		return fmt.Errorf("%s: open compressed file: %v", entry.zipEntry.Name, err)
	}
	_, err = io.Copy(output, rc)
	return err
}

type zipFormatReader struct {
	zipReader *zip.ReadCloser
	index int
}

func (rfr *zipFormatReader) Close() error {
	return nil
}

func (rfr *zipFormatReader) OpenPath(path string) error {
	var err error
	rfr.zipReader, err = zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("read: failed to create reader: %v", err)
	}
	return nil
}

// Read extracts the RAR file read from input and puts the contents
// into destination.
func (rfr *zipFormatReader) ReadEntry() (Entry, error) {
	if rfr.index >= len(rfr.zipReader.File) {
		return NilEntry, nil
	}

	f := rfr.zipReader.File[rfr.index]

	rfr.index++
	return &zipFormatEntry{f}, nil
}