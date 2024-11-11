package manga

type (
	Provider interface {
		ExtractMangaID(URL string) (string, error)
		FindManga(mangaID string) (*Manga, error)
		FindChapters(mangaID string) ([]*Chapter, error)
		FindChapter(chapterID string) (*Chapter, error)
		Extract(chapter *Chapter) ([]*Page, error)
	}

	Manga struct {
		Provider string  `json:"provider"`
		ID       string  `json:"id"`
		Title    string  `json:"title"`
		Cover    *string `json:"cover"`
		URL      string  `json:"url"`
	}

	Chapter struct {
		Provider string `json:"provider"`
		ID       string `json:"id"`
		Number   string `json:"number"`
		Title    string `json:"title"`
		Index    int    `json:"index"`
		URL      string `json:"url"`
		MangaID  string `json:"mangaId"`
	}

	DescrambleFn func(bytes []byte) ([]byte, error)

	Page struct {
		Provider     string `json:"provider"`
		Index        int    `json:"index"`
		URL          string `json:"url"`
		Filename     string `json:"filename"`
		DescrambleFn *DescrambleFn
		Headers      map[string]string `json:"headers"`
	}

	Error string
)
