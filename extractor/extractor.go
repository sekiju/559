package extractor

import (
	"fmt"
	"github.com/sekiju/mdl/extractor/cmoa"
	"github.com/sekiju/mdl/extractor/comic_walker"
	"github.com/sekiju/mdl/extractor/corocoro"
	"github.com/sekiju/mdl/extractor/ganma"
	"github.com/sekiju/mdl/extractor/storia_takeshobo"
	"github.com/sekiju/mdl/extractor/template/giga_viewer"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
)

type Factory func(cookieString *string) (manga.Extractor, error)

var domainRegistry = map[string]Factory{
	"comic-walker.com":          factorize(comic_walker.New),
	"shonenjumpplus.com":        factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-zenon.com":           factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"pocket.shonenmagazine.com": factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-gardo.com":           factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"magcomi.com":               factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"tonarinoyj.jp":             factorize(giga_viewer.New),
	"comic-ogyaaa.com":          factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-action.com":          factorize(giga_viewer.New),
	"comic-days.com":            factorize(giga_viewer.New),
	"comic-growl.com":           factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-earthstar.com":       factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"comicborder.com":           factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"comic-trail.com":           factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"kuragebunch.com":           factorize(giga_viewer.New),
	"viewer.heros-web.com":      factorize(giga_viewer.New),
	"www.sunday-webry.com":      factorizeAuthorizationOptional(giga_viewer.New, giga_viewer.NewAuthorized),
	"www.cmoa.jp":               factorizeAuthorizationRequired(cmoa.New),
	"www.corocoro.jp":           factorizeAuthorizationOptional(corocoro.New, corocoro.NewAuthorized),
	"storia.takeshobo.co.jp":    factorize(storia_takeshobo.New),
	"ganma.jp":                  factorizeAuthorizationRequiredWithError(ganma.New),
}

func factorize[T func() manga.Extractor](fn T) Factory {
	return func(cookieString *string) (manga.Extractor, error) {
		return fn(), nil
	}
}

func factorizeAuthorizationRequired[T func(string) manga.Extractor](fn T) Factory {
	return func(cookieString *string) (manga.Extractor, error) {
		if cookieString != nil {
			return fn(*cookieString), nil
		}
		return nil, manga.ErrCredentialsRequired
	}
}

func factorizeAuthorizationRequiredWithError[T func(string) (manga.Extractor, error)](fn T) Factory {
	return func(cookieString *string) (manga.Extractor, error) {
		if cookieString != nil {
			return fn(*cookieString)
		}
		return nil, manga.ErrCredentialsRequired
	}
}

func factorizeAuthorizationOptional[T func() manga.Extractor, E func(string) manga.Extractor](fn T, fnWithAuth E) Factory {
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
