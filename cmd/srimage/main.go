package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SRsawaguchi/srimage/imaging"
)

var (
	filePath = flag.String("f", "", "path to target image file.")
	destPath = flag.String("o", "", "for result image path.")
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

	fmt.Println(*filePath)
	if *destPath == "" {
		*destPath = fmt.Sprintf("out%s", filepath.Ext(*filePath))
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	dstfile, err := os.Create(*destPath)
	if err != nil {
		return err
	}
	defer dstfile.Close()

	if err := imaging.ToGrayScale(file, dstfile); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
