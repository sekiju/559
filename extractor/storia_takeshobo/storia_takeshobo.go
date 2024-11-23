package storia_takeshobo

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/extractor/speed_binb"
	"github.com/sekiju/mdl/sdk/manga"
)

type Extractor struct{}

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	//TODO implement me
	panic("implement me")
}

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	res, err := htt.New().Get(URL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      "",
		Number:  "",
		Title:   doc.Find("title").Text(),
		Index:   0,
		URL:     URL,
		MangaID: "",
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	return speed_binb.New().FindChapterPages(chapter)
}

func New() manga.Extractor {
	return new(Extractor)
}
