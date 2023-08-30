package filestore

type FileStore interface {
	Remove(string2 string) bool // Remove deletes a file by id
}
