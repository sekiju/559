package manga

type (
	ID interface{}

	Provider interface {
		ExtractMangaID(URL string) (ExtractedURL, error)
		FindManga(ID ID) (*Manga, error)
		FindChapters(ID ID) ([]*Chapter, error)
		FindChapter(ID ID) (*Chapter, error)
		ExtractPages(chapter *Chapter) ([]*Page, error)
	}

	Error string

	ExtractedURL map[string]any

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
)
