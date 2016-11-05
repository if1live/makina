package storages

type Accessor interface {
	UploadJson(data interface{}, dst string) error
	UploadBytes(data []byte, dst string) error
	UploadFile(src string, dst string) error

	Mkdir(dirname string)
}
