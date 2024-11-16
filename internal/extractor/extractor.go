package extractor

import (
	"fmt"
	"github.com/sekiju/mary/internal/config"
	"github.com/sekiju/mary/pkg/manga/comic_walker"
	"github.com/sekiju/mary/pkg/manga/giga_viewer"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
)

type ProviderFactory func(session *string) manga.Provider

var providerRegistry = map[string]ProviderFactory{
	"comic-walker.com": func(session *string) manga.Provider {
		return comic_walker.New()
	},
	"shonenjumpplus.com":        createGigaViewerProvider("shonenjumpplus.com"),
	"comic-zenon.com":           createGigaViewerProvider("comic-zenon.com"),
	"pocket.shonenmagazine.com": createGigaViewerProvider("pocket.shonenmagazine.com"),
	"comic-gardo.com":           createGigaViewerProvider("comic-gardo.com"),
	"magcomi.com":               createGigaViewerProvider("magcomi.com"),
	"tonarinoyj.jp":             createGigaViewerProviderWithoutSession("tonarinoyj.jp"),
	"comic-ogyaaa.com":          createGigaViewerProviderWithoutSession("comic-ogyaaa.com"),
	"comic-action.com":          createGigaViewerProvider("comic-action.com"),
	"comic-days.com":            createGigaViewerProvider("comic-days.com"),
	"comic-growl.com":           createGigaViewerProviderWithoutSession("comic-growl.com"),
	"comic-earthstar.com":       createGigaViewerProviderWithoutSession("comic-earthstar.com"),
	"comicborder.com":           createGigaViewerProviderWithoutSession("comicborder.com"),
	"comic-trail.com":           createGigaViewerProviderWithoutSession("comic-trail.com"),
	"kuragebunch.com":           createGigaViewerProvider("kuragebunch.com"),
	"viewer.heros-web.com":      createGigaViewerProvider("viewer.heros-web.com"),
	"www.sunday-webry.com":      createGigaViewerProviderWithoutSession("www.sunday-webry.com"),
}

func createGigaViewerProvider(hostname string) ProviderFactory {
	return func(session *string) manga.Provider {
		if session != nil {
			return giga_viewer.NewWithSession(hostname, *session)
		}
		return giga_viewer.New(hostname)
	}
}

func createGigaViewerProviderWithoutSession(hostname string) ProviderFactory {
	return func(session *string) manga.Provider {
		return giga_viewer.New(hostname)
	}
}

func CreateProvider(cfg *config.Config, args *config.Arguments, hostname string) (manga.Provider, error) {
	factory, exists := providerRegistry[hostname]
	if !exists {
		return nil, fmt.Errorf("unsupported provider for hostname: %s", hostname)
	}

	session := getSession(cfg, args, hostname)
	return factory(session), nil
}

func getSession(cfg *config.Config, args *config.Arguments, hostname string) *string {
	if args.Session != "" {
		return &args.Session
	}
	if site, exists := cfg.Sites[hostname]; exists && site.Session != nil {
		return site.Session
	}
	return nil
}
