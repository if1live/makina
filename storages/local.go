package storages

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path"

	raven "github.com/getsentry/raven-go"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
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
	dir := path.Dir(filepath)
	c.Mkdir(dir)

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

func (c *Local) UploadYaml(data interface{}, dst string) error {
	d, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	filepath := path.Join(c.RootPath, dst)
	dir := path.Dir(filepath)
	c.Mkdir(dir)

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(d)
	return nil
}

func (c *Local) UploadBytes(data []byte, dst string) error {
	filepath := path.Join(c.RootPath, dst)
	dir := path.Dir(filepath)
	c.Mkdir(dir)

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
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
