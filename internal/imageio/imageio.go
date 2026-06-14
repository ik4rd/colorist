package imageio

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const jpegQuality = 92

const (
	FormatPNG  = "png"
	FormatJPEG = "jpeg"
)

func Load(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	return img, nil
}

func Save(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	err = Encode(f, filepath.Ext(path), img)
	if cerr := f.Close(); err == nil {
		err = cerr
	}

	return err
}

func Encode(w io.Writer, format string, img image.Image) error {
	switch strings.ToLower(strings.TrimPrefix(format, ".")) {
	case "jpg", "jpeg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: jpegQuality})
	case "png", "":
		return png.Encode(w, img)
	default:
		return fmt.Errorf("imageio: unsupported format %q (use png or jpg)", format)
	}
}
