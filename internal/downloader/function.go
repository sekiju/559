package downloader

import (
	"bytes"
	"fmt"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

func pReader(page *manga.Page) (io.Reader, error) {
	res, err := htt.New().SetHeaders(page.Headers).Get(page.URL)
	if err != nil {
		return nil, err
	}

	if page.Decode == nil {
		return res.Body, nil
	}

	b, err := res.Bytes()
	if err != nil {
		return nil, err
	}

	b, err = page.Decode(b)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func Bytes(dir string, page *manga.Page) error {
	r, err := pReader(page)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, page.Filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	_, err = file.ReadFrom(r)
	if err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func WithEncode(dir string, format config.OutputFileFormat, page *manga.Page) error {
	r, err := pReader(page)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	var fileExt string
	switch format {
	case config.JpegOutputFormat:
		fileExt = "jpg"
	case config.PngOutputFormat:
		fileExt = "png"
	case config.AvifOutputFormat:
		fileExt = "avif"
	case config.WebpOutputFormat:
		fileExt = "webp"
	}

	file, err := os.Create(filepath.Join(dir, page.Filename[:1+len(page.Filename)-len(filepath.Ext(page.Filename))]+fileExt))
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	switch format {
	case config.JpegOutputFormat:
		if err = jpeg.Encode(file, img, nil); err != nil {
			return err
		}
	case config.PngOutputFormat:
		if err = png.Encode(file, img); err != nil {
			return err
		}
	case config.AvifOutputFormat:
		if err = avif.Encode(file, img); err != nil {
			return err
		}
	case config.WebpOutputFormat:
		if err = webp.Encode(file, img); err != nil {
			return err
		}
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}
