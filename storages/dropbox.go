package storages

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

type Dropbox struct {
	RootPath string

	client *dropy.Client
	token  string
}

func NewDropbox(rootpath string, token string) *Dropbox {
	client := dropy.New(dropbox.New(dropbox.NewConfig(token)))
	return &Dropbox{
		RootPath: rootpath,
		client:   client,
		token:    token,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (c *Dropbox) UploadJson(data interface{}, dst string) error {
	b, err := json.Marshal(data)

	check(err)
	var jsonOut bytes.Buffer
	json.Indent(&jsonOut, b, "", "  ")

	r := bytes.NewReader(jsonOut.Bytes())
	uploadFilePath := path.Join(c.RootPath, dst)
	dir := path.Dir(uploadFilePath)
	c.Mkdir(dir)
	e := c.client.Upload(uploadFilePath, r)
	return e
}

func (c *Dropbox) UploadYaml(data interface{}, dst string) error {
	d, err := yaml.Marshal(data)
	check(err)

	r := bytes.NewReader(d)
	uploadFilePath := path.Join(c.RootPath, dst)
	dir := path.Dir(uploadFilePath)
	c.Mkdir(dir)
	e := c.client.Upload(uploadFilePath, r)
	return e
}

func (c *Dropbox) UploadBytes(data []byte, dst string) error {
	r := bytes.NewReader(data)
	uploadFilePath := path.Join(c.RootPath, dst)
	dir := path.Dir(uploadFilePath)
	c.Mkdir(dir)
	err := c.client.Upload(uploadFilePath, r)
	return err
}

func (c *Dropbox) UploadFile(src string, dst string) error {
	uploadFilePath := path.Join(c.RootPath, dst)
	dir := path.Dir(uploadFilePath)
	c.Mkdir(dir)
	file, _ := os.Open(src)
	r := bufio.NewReader(file)
	err := c.client.Upload(uploadFilePath, r)
	return err
}

func (c *Dropbox) Mkdir(dirname string) {
	client := c.client
	_, err := client.Stat(dirname)
	if err != nil {
		client.Mkdir(dirname)
	}
}
