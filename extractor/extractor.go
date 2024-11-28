package extractor

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/extractor/cmoa"
	"github.com/sekiju/mdl/extractor/comic_walker"
	"github.com/sekiju/mdl/extractor/corocoro"
	"github.com/sekiju/mdl/extractor/ganma"
	"github.com/sekiju/mdl/extractor/storia_takeshobo"
	"github.com/sekiju/mdl/extractor/template/giga_viewer"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
)

type Factory func(cookie *string) (manga.Extractor, error)

var domainRegistry = map[string]Factory{
	"comic-walker.com":          fz(comic_walker.New),
	"shonenjumpplus.com":        fz(giga_viewer.New),
	"comic-zenon.com":           fz(giga_viewer.New),
	"pocket.shonenmagazine.com": fz(giga_viewer.New),
	"comic-gardo.com":           fz(giga_viewer.New),
	"magcomi.com":               fz(giga_viewer.New),
	"tonarinoyj.jp":             fz(giga_viewer.New),
	"comic-ogyaaa.com":          fz(giga_viewer.New),
	"comic-action.com":          fz(giga_viewer.New),
	"comic-days.com":            fz(giga_viewer.New),
	"comic-growl.com":           fz(giga_viewer.New),
	"comic-earthstar.com":       fz(giga_viewer.New),
	"comicborder.com":           fz(giga_viewer.New),
	"comic-trail.com":           fz(giga_viewer.New),
	"kuragebunch.com":           fz(giga_viewer.New),
	"viewer.heros-web.com":      fz(giga_viewer.New),
	"www.sunday-webry.com":      fz(giga_viewer.New),
	"www.cmoa.jp":               fz(cmoa.New),
	"www.corocoro.jp":           fz(corocoro.New),
	"storia.takeshobo.co.jp":    fz(storia_takeshobo.New),
	"ganma.jp":                  fz(ganma.New),
}

// fz is a generic helper function to create a Factory for manga.Extractor
func fz[T func() (manga.Extractor, error)](fn T) Factory {
	return func(cookie *string) (manga.Extractor, error) {
		ext, err := fn()
		if err != nil {
			return nil, err
		}

		if cookie != nil {
			ext.SetSettings(manga.Settings{Cookie: cookie})
		} else {
			cookieGenerator, ok := ext.(manga.GenerateCookieFeature)
			if ok {
				generatedCookie, err := cookieGenerator.GenerateCookie()
				if err != nil {
					return nil, err
				}

				log.Info().Msgf("Cookie generated >>> %s", generatedCookie)

				ext.SetSettings(manga.Settings{Cookie: &generatedCookie})
			}
		}

		return ext, nil
	}
}

func getSession(cfg *config.Config, hostname string) *string {
	if cfg.PrimaryCookie != nil {
		return cfg.PrimaryCookie
	}
	if site, exists := cfg.Sites[hostname]; exists && site.Cookie != nil {
		return site.Cookie
	}
	return nil
}

func NewExtractor(cfg *config.Config, hostname string) (manga.Extractor, error) {
	factory, exists := domainRegistry[hostname]
	if !exists {
		return nil, fmt.Errorf("unsupported website: %s", hostname)
	}
	return factory(getSession(cfg, hostname))
}
