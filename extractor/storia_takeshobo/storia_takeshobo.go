package storia_takeshobo

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/extractor/template/speed_binb"
	"github.com/sekiju/mdl/sdk/manga"
	"regexp"
)

type Extractor struct{}

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	//TODO implement me
	panic("implement me")
}

var re = regexp.MustCompile(`https://storia.takeshobo.co.jp/_files/([a-zA-Z0-9_]*)/(\d*)`)

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	res, err := htt.New().Get(URL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      matches[2],
		Number:  "",
		Title:   doc.Find("title").Text(),
		Index:   0,
		URL:     URL,
		MangaID: matches[1],
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	return speed_binb.New().FindChapterPages(chapter)
}

func New() manga.Extractor {
	return new(Extractor)
}
