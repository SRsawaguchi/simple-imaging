package storage

import (
	"io"
	"io/ioutil"
	"os"
)

type Downloader interface {
	Download(key string, dest io.WriterAt) error
}

type Uploader interface {
	Upload(key string, src io.Reader) error
	GenerateUploadKey(fileName string) string
}

type TempFileManager interface {
	MakeTempDir(dir, prefix string) (string, error)
	MakeTempFile(dir, pattern string) (*os.File, error)
}

type TempFileManagerImp struct{}

func (t TempFileManagerImp) MakeTempDir(dir, prefix string) (string, error) {
	dir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func (t TempFileManagerImp) MakeTempFile(dir, pattern string) (*os.File, error) {
	file, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, err
	}

	return file, nil
}
