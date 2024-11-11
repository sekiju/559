package giga_viewer

import (
	"github.com/sekiju/mary/pkg/manga/internal/util"
	"github.com/sekiju/mary/pkg/sdk/extractor/manga"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	provider := New("shonenjumpplus.com")

	t.Run("ExtractMangaID", func(t *testing.T) {
		_, err := provider.ExtractMangaID("https://shonenjumpplus.com/episode/17106371892806261346")
		assert.Equal(t, manga.ErrURLeqID, err)

		_, err = provider.ExtractMangaID("https://shonenjumpplus.com/magazine/4856001361589090626")
		assert.Equal(t, manga.ErrURLeqID, err)
	})

	t.Run("FindManga", func(t *testing.T) {
		data, err := provider.FindManga("https://shonenjumpplus.com/episode/17106371892806261346")
		assert.NoError(t, err)
		assert.Equal(t, "魔都精兵のスレイブ", data.Title)
	})

	t.Run("FindChapters", func(t *testing.T) {
		episodes, err := provider.FindChapters("https://shonenjumpplus.com/episode/17106371892806261346")
		assert.NoError(t, err)
		assert.NotEmpty(t, episodes)
		assert.Equal(t, "17106371892806261346", episodes[0].ID)
	})

	t.Run("FindChapter", func(t *testing.T) {
		chapter, err := provider.FindChapter("https://shonenjumpplus.com/episode/17106371892806261346")
		assert.NoError(t, err)
		assert.Equal(t, "17106371892806261346", chapter.ID)

		t.Run("ExtractEpisode", func(t *testing.T) {
			pages, err := provider.ExtractPages(chapter)
			assert.NoError(t, err)
			assert.NotEmpty(t, pages)
			util.AssertImage(t, "https://stg.yandere.ovh/test_providers/giga_viewer__sjp__episode%2417106371892806261346.jpg", pages[0])
		})
	})
}
