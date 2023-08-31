package filestore

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"os"
)

const randIDLength = 5

var errUniqueIDNotGenerated = errors.New("unique id does not exists after tried 50 times")

type FileStore interface {
	Remove(string2 string) bool            // Remove 通过文件id删除文件
	List() map[string]string               // List 返回所有文件id的原始名称
	New() (*os.File, error)                // 创建一个临时文件到文件存储，可以通过添加来保存它
	Add(name, path string) (string, error) // Add 创建一个带有存储路径的文件，返回id
}

func generateID() (string, error) {
	b := make([]byte, randIDLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(b), nil
}
