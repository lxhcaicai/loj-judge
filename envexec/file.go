package envexec

import (
	"fmt"
	"io"
	"os"
)

type File interface {
	isFile()
}

// FileReader 表示在执行之前可以完全读取的文件输入
type FileReader struct {
	Reader io.Reader
	Stream bool
}

func (*FileReader) isFile() {}

type FileInput struct {
	Path string
}

func (*FileInput) isFile() {}

// FileOpened 表示已打开的文件
type FileOpened struct {
	File *os.File
}

func (*FileOpened) isFile() {}

func NewFileInput(p string) File {
	return &FileInput{Path: p}
}

func FileToReader(f File) (io.ReadCloser, error) {
	switch f := f.(type) {
	case *FileOpened:
		return f.File, nil
	case *FileReader:
		return io.NopCloser(f.Reader), nil
	case *FileInput:
		file, err := os.Open(f.Path)
		if err != nil {
			return nil, err
		}
		return file, nil
	default:
		return nil, fmt.Errorf("file cannot open as reader %v", f)
	}
}
