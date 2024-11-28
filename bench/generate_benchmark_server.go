package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/extractor/util"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	chapterURL = "https://comic-ogyaaa.com/episode/2550912964910100594"
)

/* update hosts
127.0.0.1 comic-ogyaaa.com
127.0.0.1 cdn-img.comic-ogyaaa.com
*/

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func run() error {
	res, err := htt.New().SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1").Get(chapterURL)
	if err != nil {
		return err
	}

	html, err := res.Text()
	if err != nil {
		return err
	}

	episode, err := util.ExtractJSONFromHTML[episodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
	if err != nil {
		return err
	}

	for _, page := range episode.ReadableProduct.PageStructure.Pages {
		if page.Type != "main" {
			continue
		}

		parsedURL, _ := url.Parse(page.Src)
		filePath := filepath.Join(".", parsedURL.Path)

		if _, err = os.Stat(filePath); err == nil {
			continue
		}

		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		res, err = htt.New().Get(page.Src)
		if err != nil {
			return err
		}

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		_, err = file.ReadFrom(res.Body)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		if err = file.Close(); err != nil {
			return err
		}
	}

	parsedURL, _ := url.Parse(chapterURL)
	filePath := filepath.Join(".", parsedURL.Path)

	if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	html = strings.ReplaceAll(html, "https://", "http://")

	if err = os.WriteFile(filePath, []byte(html), os.ModePerm); err != nil {
		return err
	}

	app := fiber.New()

	app.Static("/", ".")

	return app.Listen(":80")
}

type episodeResult struct {
	ReadableProduct struct {
		FinishReadingNotificationUri interface{} `json:"finishReadingNotificationUri"`
		HasPurchased                 bool        `json:"hasPurchased"`
		Id                           string      `json:"id"`
		IsPublic                     bool        `json:"isPublic"`
		NextReadableProductUri       *string     `json:"nextReadableProductUri"`
		Number                       int         `json:"number"`
		PageStructure                *struct {
			Pages []episodeResultPage `json:"pages"`
		} `json:"pageStructure"`
		Permalink              string  `json:"permalink"`
		PrevReadableProductUri *string `json:"prevReadableProductUri"`
		Title                  string  `json:"title"`
	} `json:"readableProduct"`
}

type episodeResultPage struct {
	Type string `json:"type"`
	Src  string `json:"src,omitempty"`
}
