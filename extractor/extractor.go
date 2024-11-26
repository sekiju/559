package extractor

import (
	"fmt"
	"github.com/sekiju/mdl/extractor/cmoa"
	"github.com/sekiju/mdl/extractor/comic_walker"
	"github.com/sekiju/mdl/extractor/corocoro"
	"github.com/sekiju/mdl/extractor/giga_viewer"
	"github.com/sekiju/mdl/extractor/storia_takeshobo"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
)

type Factory func(cookieString *string) (manga.Extractor, error)

var domainRegistry = map[string]Factory{
	"comic-walker.com":          factorize(comic_walker.New),
	"shonenjumpplus.com":        factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-zenon.com":           factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"pocket.shonenmagazine.com": factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-gardo.com":           factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"magcomi.com":               factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"tonarinoyj.jp":             factorize(giga_viewer.New),
	"comic-ogyaaa.com":          factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-action.com":          factorize(giga_viewer.New),
	"comic-days.com":            factorize(giga_viewer.New),
	"comic-growl.com":           factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-earthstar.com":       factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"comicborder.com":           factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-trail.com":           factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"kuragebunch.com":           factorize(giga_viewer.New),
	"viewer.heros-web.com":      factorize(giga_viewer.New),
	"www.sunday-webry.com":      factorizeAuthorized(giga_viewer.New, giga_viewer.NewAuthorized),
	"www.cmoa.jp": func(cookieString *string) (manga.Extractor, error) {
		if cookieString != nil {
			return cmoa.New(*cookieString), nil
		}
		return nil, manga.ErrCredentialsRequired
	},
	"www.corocoro.jp":        factorizeAuthorized(corocoro.New, corocoro.NewAuthorized),
	"storia.takeshobo.co.jp": factorize(storia_takeshobo.New),
}

func factorize[T func() manga.Extractor](fn T) Factory {
	return func(cookieString *string) (manga.Extractor, error) {
		return fn(), nil
	}
}

// factorizeAuthorized used for create Factory with optional authorization
func factorizeAuthorized[T func() manga.Extractor, E func(string) manga.Extractor](fn T, fnWithAuth E) Factory {
	return func(cookieString *string) (manga.Extractor, error) {
		if cookieString != nil {
			return fnWithAuth(*cookieString), nil
		}

		return fn(), nil
	}
}

func getSession(cfg *config.Config, hostname string) *string {
	if cfg.PrimaryCookie != nil {
		return cfg.PrimaryCookie
	}

	if site, exists := cfg.Sites[hostname]; exists && site.CookieString != nil {
		return site.CookieString
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
