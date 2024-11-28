package main

import (
	"github.com/sekiju/mdl/downloader"
	"testing"
)

func Benchmark(b *testing.B) {
	chapterURLs := []string{"http://comic-ogyaaa.com/episode/2550912964910100594"}

	b.Run("Raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			loader := downloader.NewDownloader()

			for _, chapterURL := range chapterURLs {
				loader.Queue(chapterURL)
			}

			loader.Stop()
		}
	})
}
