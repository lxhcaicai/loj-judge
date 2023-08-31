package filestore

type FileStore interface {
	Remove(string2 string) bool // Remove 通过文件id删除文件
	List() map[string]string    // List 返回所有文件id的原始名称
}
