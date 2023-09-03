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

// FileWriter 表示将通过管道从exec输出的管道输出
type FileWriter struct {
	Writer io.Writer
	Limit  Size
}

func (*FileWriter) isFile() {}

func NewFileInput(p string) File {
	return &FileInput{Path: p}
}

// FileCollector 表示将通过管道收集的管道输出
type FileCollector struct {
	Name  string
	Limit Size
	Pipe  bool
}

func (f FileCollector) isFile() {}

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

// NewFileReader creates File input which can be fully read before exec
// or piped into exec
func NewFileReader(r io.Reader, s bool) File {
	return &FileReader{
		Reader: r,
		Stream: s,
	}
}

// NewFileCollector 创建将通过管道收集的文件输出
func NewFileCollector(name string, limit Size, pipe bool) File {
	return &FileCollector{
		Name:  name,
		Limit: limit,
		Pipe:  pipe,
	}
}

// ReaderTTY will be asserts when File Reader is provided and TTY is enabled
// and then TTY will be called with pty file
type ReaderTTY interface {
	TTY(file *os.File)
}
