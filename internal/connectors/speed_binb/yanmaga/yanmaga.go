package yanmaga

import (
	"github.com/sekiju/rq"
	"mary/internal/config"
	"mary/internal/connectors/speed_binb"
	"mary/internal/static"
	"mary/internal/utils"
	"net/url"
)

type Yanmaga struct {
	domain string
	binb   *speed_binb.SpeedBinb
}

func (c *Yanmaga) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusOptional,
		ChapterListAvailable: true,
	}
}

func (c *Yanmaga) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *Yanmaga) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *Yanmaga) Chapter(uri url.URL) (*static.Chapter, error) {
	document, err := utils.Document(uri.String(), c.withCookies())
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    uri,
		Title: document.Find("title").Text(),
		Error: nil,
	}, nil
}

func (c *Yanmaga) Pages(chapterID any, imageChan chan<- static.Image) error {
	opts := c.withCookies()
	return c.binb.Pages(chapterID.(url.URL), imageChan, &opts)
}

func (c *Yanmaga) withCookies() rq.OptsFn {
	connectorConfig, exists := config.Config.Sites[c.domain]
	return func(cf *rq.Opts) {
		if exists {
			cf.Headers["cookie"] = connectorConfig.Session
		}
	}
}

func New() *Yanmaga {
	domain := "yanmaga.jp"
	return &Yanmaga{
		domain: domain,
		binb:   speed_binb.New(domain),
	}
}
