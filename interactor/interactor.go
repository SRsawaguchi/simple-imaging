package interactor

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/SRsawaguchi/srimage/imaging"
	"github.com/SRsawaguchi/srimage/storage"
)

type Interactor interface {
	Execute() (*Result, error)
	Clean() error
}

type Result struct {
	ResultImageLocalPath string
	UploadedImageKey     string
}

type interactorImp struct {
	Downloader      storage.Downloader
	Uploader        storage.Uploader
	TempFileManager storage.TempFileManager
	workDirRoot     string
	workDir         string
	srcFileKey      string
}

func NewInteractor(
	downloader storage.Downloader,
	uploader storage.Uploader,
	workDirRoot string,
	srcFileKey string,
) *interactorImp {
	return &interactorImp{
		Downloader:      downloader,
		Uploader:        uploader,
		workDirRoot:     workDirRoot,
		srcFileKey:      srcFileKey,
		TempFileManager: storage.TempFileManagerImp{},
	}
}

func (i interactorImp) openSrcFile(workDir string) (*os.File, error) {
	if i.Downloader == nil {
		srcFile, err := os.Open(i.srcFileKey)
		if err != nil {
			return nil, err
		}

		return srcFile, nil
	}

	ext := path.Ext(i.srcFileKey)
	srcFile, err := i.TempFileManager.MakeTempFile(workDir, "srimage.*"+ext)
	if err != nil {
		return nil, err
	}
	if err = i.Downloader.Download(i.srcFileKey, srcFile); err != nil {
		return nil, err
	}

	return srcFile, nil
}

func (i interactorImp) openDestFile(workDir string) (*os.File, error) {
	ext := path.Ext(i.srcFileKey)
	destFile, err := i.TempFileManager.MakeTempFile(workDir, "srimage.result.*"+ext)
	if err != nil {
		return nil, err
	}

	return destFile, err
}

func (i interactorImp) upload(file io.Reader) (string, error) {
	if i.Uploader == nil {
		return "", nil
	}

	key := i.Uploader.GenerateUploadKey(filepath.Base(i.srcFileKey))
	if err := i.Uploader.Upload(key, file); err != nil {
		return "", err
	}

	return key, nil
}

func (i *interactorImp) Execute() (*Result, error) {
	workDir, err := i.TempFileManager.MakeTempDir(i.workDirRoot, "")
	if err != nil {
		return nil, err
	}
	i.workDir = workDir

	srcFile, err := i.openSrcFile(workDir)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	destFile, err := i.openDestFile(workDir)
	if err != nil {
		return nil, err
	}
	defer destFile.Close()

	if err = imaging.ToGrayScale(srcFile, destFile); err != nil {
		return nil, err
	}
	destFile.Close()

	// for rewind file pointer
	destFile, err = os.Open(destFile.Name())
	if err != nil {
		return nil, err
	}
	defer destFile.Close()

	key, err := i.upload(destFile)
	if err != nil {
		return nil, err
	}

	result := &Result{
		ResultImageLocalPath: destFile.Name(),
		UploadedImageKey:     key,
	}
	return result, nil
}

func (i interactorImp) Clean() error {
	if i.workDir == "" {
		return nil
	}

	if err := os.RemoveAll(i.workDir); err != nil {
		return err
	}
	return nil
}
