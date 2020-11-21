package file_store

type FileStoreType interface {
	Setup() error
	UpLoad(yourObjectName string, localFile string) error
}