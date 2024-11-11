package download

import (
	"fmt"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/sekiju/rq"
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
