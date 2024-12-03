package downloader

import (
	"bytes"
	"fmt"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/config"
	"github.com/sekiju/mdl/sdk/manga"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

func getReader(page *manga.Page) (io.Reader, error) {
	res, err := htt.New().SetHeaders(page.Headers).Get(page.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download page: %w", err)
	}

	if page.Decode == nil {
		return res.Body, nil
	}

	b, err := res.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	b, err = page.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode page: %w", err)
	}

	return bytes.NewReader(b), nil
}

func saveFile(dir, filename string, r io.Reader) error {
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil

}

func generateFilename(originalName string, format config.OutputFileFormat) string {
	return originalName[:1+len(originalName)-len(filepath.Ext(originalName))] + string(format)
}

func saveEncodedImage(dir, filename string, format config.OutputFileFormat, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	filename = generateFilename(filename, format)
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch format {
	case config.JpegOutputFormat:
		return jpeg.Encode(file, img, nil)
	case config.PngOutputFormat:
		return png.Encode(file, img)
	case config.AvifOutputFormat:
		return avif.Encode(file, img)
	case config.WebpOutputFormat:
		return webp.Encode(file, img)
	default:
		return fmt.Errorf("unsupported output format: %v", format)
	}
}
