package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SRsawaguchi/srimage/imaging"
	"github.com/SRsawaguchi/srimage/interactor"
)

var (
	filePath = flag.String("f", "", "path to target image file.")
	destPath = flag.String("o", "", "for result image path.")
	workdir  = flag.String("w", "", "work directory")
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func run() error {
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please specify a file.")
		return nil
	}

	if !imaging.IsImageFile(*filePath) {
		fmt.Println("This is not a image file.")
		return nil
	}

	if *destPath == "" {
		*destPath = fmt.Sprintf("out%s", filepath.Ext(*filePath))
	}

	if *workdir == "" {
		*workdir = os.Getenv("SRIMAGE_WORKDIR")
		if *workdir == "" {
			*workdir = "./workspace"
			if err := os.Mkdir(*workdir, 0777); err != nil {
				return err
			}
		}
	}

	interactor := interactor.NewInteractor(
		nil,
		nil,
		*workdir,
		*filePath,
	)

	result, err := interactor.Execute()
	if err != nil {
		return err
	}

	if err = os.Rename(result.ResultImageLocalPath, *destPath); err != nil {
		return err
	}

	interactor.Clean()

	return nil
}
