package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

type IFile interface {
	Name() string
	Value() []byte
	Reader() io.Reader
	Writer() io.Writer
	File() (*os.File, error)
}

type File struct {
	name  string
	path  string
	value []byte
}

func NewFile(name string, value []byte) IFile {
	path := fmt.Sprintf("/tmp/%s", uuid.New().String())

	return &File{
		name:  name,
		value: value,
		path:  path,
	}
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Value() []byte {
	return f.value
}

func (f *File) Reader() io.Reader {
	return bytes.NewReader(f.Value())
}

func (f *File) Writer() io.Writer {
	return bytes.NewBuffer(f.Value())
}

func (f *File) File() (*os.File, error) {
	tmp, err := os.Create(f.path)
	if err != nil {
		return nil, err
	}

	_, err = tmp.Write(f.value)
	if err != nil {
		return nil, err
	}

	err = tmp.Close()
	if err != nil {
		return nil, err
	}

	return os.Open(f.path)
}
