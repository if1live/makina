package storages

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path"

	"github.com/kardianos/osext"
)

type Local struct {
	RootPath string
}

func NewLocal() *Local {
	rootpath, _ := osext.ExecutableFolder()
	return &Local{
		RootPath: rootpath,
	}
}

func (c *Local) UploadJson(data interface{}, dst string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	filepath := path.Join(c.RootPath, dst)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	out.WriteTo(w)

	w.Flush()
	return nil
}

func (c *Local) UploadBytes(data []byte, dst string) error {
	filepath := path.Join(c.RootPath, dst)
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	file.Write(data)
	defer file.Close()
	return err
}

func (c *Local) UploadFile(src string, dst string) error {
	panic("not implemented")
}

func (c *Local) Mkdir(dirname string) {
	os.MkdirAll(dirname, 0644)
}
