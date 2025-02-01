package bilibili

import (
	"github.com/goccy/go-json"
	"github.com/sekiju/htt"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/sdk/manga"
	"regexp"
	"strconv"
)

type Extractor struct {
	settings *manga.Settings
}

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	//TODO implement me
	panic("implement me")
}

var re = regexp.MustCompile(`https://manga.bilibili.com/([a-zA-Z0-9_]*)/(\d*)`)

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	res, err := htt.New().Body(map[string]interface{}{"id": matches[2]}).
		Post("https://manga.bilibili.com/twirp/comic.v1.Comic/GetEpisode?device=pc&platform=web&nov=25")
	if err != nil {
		return nil, err
	}

	var episode getEpisodeResult
	if err = res.JSON(&episode); err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      matches[2],
		Number:  episode.Data.ShortTitle,
		Title:   episode.Data.Title,
		Index:   0,
		URL:     URL,
		MangaID: strconv.Itoa(episode.Data.ComicId),
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	episodeID, _ := strconv.Atoi(chapter.ID)

	res, err := htt.New().Body(map[string]interface{}{"ep_id": episodeID}).
		Post("https://manga.bilibili.com/twirp/comic.v1.Comic/GetImageIndex?device=pc&platform=web&nov=25")
	if err != nil {
		return nil, err
	}

	var imageIndex getImageIndexResult
	if err = res.JSON(&imageIndex); err != nil {
		return nil, err
	}

	imageURLs := make([]string, len(imageIndex.Data.Images))
	for i, image := range imageIndex.Data.Images {
		imageURLs[i] = image.Path
	}

	serializedImageURLs, err := json.Marshal(imageURLs)
	if err != nil {
		return nil, err
	}

	res, err = htt.New().
		Body(map[string]interface{}{
			// todo: how to generate m1?
			"m1":   "___",
			"urls": string(serializedImageURLs),
		}).
		Post("https://manga.bilibili.com/twirp/comic.v1.Comic/ImageToken?device=pc&platform=web&nov=25")
	if err != nil {
		return nil, err
	}

	var imageToken imageTokenResult
	if err = res.JSON(&imageToken); err != nil {
		return nil, err
	}

	pages := make([]*manga.Page, len(imageToken.Data))
	padRenamer := renamer.New(len(imageToken.Data))

	for i, page := range imageToken.Data {
		pages[i] = &manga.Page{
			Index:    uint(i),
			URL:      page.CompleteUrl,
			Filename: padRenamer.Name(i, ".jpg"),
		}
	}

	return pages, nil
}

func (e *Extractor) SetSettings(settings manga.Settings) {
	e.settings = &settings
}

func New() (manga.Extractor, error) {
	return &Extractor{settings: &manga.Settings{}}, nil
}
