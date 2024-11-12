package download

import (
	bt "bytes"
	"fmt"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/sekiju/mary/internal/config"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/sekiju/rq"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func Bytes(dir string, page *manga.Page) error {
	res, err := rq.New(rq.SetHeaders(page.Headers)).Get(page.URL)
	if err != nil {
		return err
	}

	bytes, err := res.Bytes()
	if page.DescrambleFn != nil {
		bytes, err = (*page.DescrambleFn)(bytes)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filepath.Join(dir, page.Filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func WithEncode(dir string, format config.OutputFormat, page *manga.Page) error {
	res, err := rq.New(rq.SetHeaders(page.Headers)).Get(page.URL)
	if err != nil {
		return err
	}

	bytes, err := res.Bytes()
	if page.DescrambleFn != nil {
		bytes, err = (*page.DescrambleFn)(bytes)
		if err != nil {
			return err
		}
	}

	img, _, err := image.Decode(bt.NewReader(bytes))
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
