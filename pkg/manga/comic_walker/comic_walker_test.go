package comic_walker

import (
	"github.com/sekiju/mary/pkg/manga/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	provider := NewProvider()

	t.Run("ExtractMangaID", func(t *testing.T) {
		mangaWithChapterID, err := provider.ExtractMangaID("https://comic-walker.com/detail/KC_005558_S/episodes/KC_0055580000200011_E?episodeType=first")
		assert.NoError(t, err)
		assert.Equal(t, "KC_005558_S$KC_0055580000200011_E", mangaWithChapterID)

		mangaID, err := provider.ExtractMangaID("https://comic-walker.com/detail/KC_005558_S?episodeType=first")
		assert.NoError(t, err)
		assert.Equal(t, "KC_005558_S$", mangaID)
	})

	t.Run("FindManga", func(t *testing.T) {
		data, err := provider.FindManga("KC_005558_S$")
		assert.NoError(t, err)
		assert.Equal(t, "忍者の騎士", data.Title)
	})

	t.Run("FindChapters", func(t *testing.T) {
		episodes, err := provider.FindChapters("KC_005558_S$")
		assert.NoError(t, err)
		assert.NotEmpty(t, episodes)
		assert.Equal(t, "018f84b1-1d0b-7557-b1e2-7ec22323c494", episodes[0].ID)
	})

	t.Run("FindChapter", func(t *testing.T) {
		chapter, err := provider.FindChapter("KC_005558_S$KC_0055580000200011_E")
		assert.NoError(t, err)
		assert.Equal(t, "018f84b1-1d0b-7557-b1e2-7ec22323c494", chapter.ID)

		t.Run("ExtractEpisode", func(t *testing.T) {
			pages, err := provider.Extract(chapter)
			assert.NoError(t, err)
			assert.NotEmpty(t, pages)
			util.AssertImage(t, "https://stg.yandere.ovh/test_providers/comic_walker__KC_005558_S%24KC_0055580000200011_E.webp", pages[0])
		})
	})
}
