package imaging

import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"strings"
)

func hasExtension(fileName string, extensions ...string) bool {
	fName := "." + fileName
	for _, ext := range extensions {
		if strings.HasSuffix(fName, ext) || strings.HasSuffix(fName, strings.ToUpper(ext)) {
			return true
		}
	}
	return false
}

func IsImageFile(fileName string) bool {
	return hasExtension(fileName, "png", "jpg", "jpeg")
}

func ToGrayScale(file io.Reader, dstfile io.Writer) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	dest := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(img.At(x, y))
			gray, _ := c.(color.Gray16)
			dest.Set(x, y, gray)
		}
	}

	return jpeg.Encode(dstfile, dest, nil)
}
