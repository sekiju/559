package manga

import "errors"

var (
	ErrMethodUnimplemented = errors.New("method unimplemented")
	ErrMangaNotFound       = errors.New("manga not found")
	ErrChapterNotFound     = errors.New("chapter not found")
	ErrPaidChapter         = errors.New("chapter is paid")
	ErrInvalidURLFormat    = errors.New("invalid URL format")
	ErrInvalidChapterURL   = errors.New("cannot extract chapterID from this url")
	ErrCredentialsRequired = errors.New("credentials required: set site.'domain'.cookie field in config.hcl")
)
