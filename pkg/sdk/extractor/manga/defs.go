package manga

const (
	ErrMethodUnimplemented Error = "method unimplemented"
	ErrMangaNotFound       Error = "manga not found"
	ErrChapterNotFound     Error = "chapter not found"
	ErrPaidChapter         Error = "chapter is paid"
	ErrInvalidURLFormat    Error = "invalid URL format"
	ErrURLIsID             Error = "URL is ID"
)

type Error string

func (e Error) Error() string {
	return string(e)
}
