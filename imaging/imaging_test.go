package imaging

import (
	"fmt"
	"os"
	"testing"
)

const testImage = "testdata/lena.jpg"

func TestIsImaging(t *testing.T) {
	extensions := []string{"jpg", "png", "jpeg", "JPG"}
	for _, ext := range extensions {
		if !IsImageFile(fmt.Sprintf("image.%s", ext)) {
			t.Fatal(fmt.Sprintf("%s is a extension of image file.", ext))
		}
	}
}

type MockWriter struct {
	fileSize int
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	l := len(p)
	m.fileSize += l
	return l, nil
}

func TestToGrayScale(t *testing.T) {
	file, err := os.Open(testImage)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer file.Close()

	writer := &MockWriter{}
	err = ToGrayScale(file, writer)

	if err != nil {
		t.Fatal(err.Error())
	}

	if writer.fileSize == 0 {
		t.Fatal("There was not output data.")
	}
}
