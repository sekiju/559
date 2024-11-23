package manga

import "errors"

var (
	ErrMethodUnimplemented   = errors.New("method unimplemented")
	ErrMangaNotFound         = errors.New("manga not found")
	ErrChapterNotFound       = errors.New("chapter not found")
	ErrPaidChapter           = errors.New("chapter is paid")
	ErrInvalidURLFormat      = errors.New("invalid URL format")
	ErrStringAsIDUnsupported = errors.New("string as ID is unsupported, please use ExtractMangaID and pass it to methods")
	ErrInvalidID             = errors.New("invalid ID")
	ErrInvalidChapterURL     = errors.New("cannot extract chapterID from this url")
)
