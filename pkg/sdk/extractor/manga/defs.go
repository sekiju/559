package manga

const (
	ErrMethodUnimplemented   Error = "method unimplemented"
	ErrMangaNotFound         Error = "manga not found"
	ErrChapterNotFound       Error = "chapter not found"
	ErrPaidChapter           Error = "chapter is paid"
	ErrInvalidURLFormat      Error = "invalid URL format"
	ErrStringAsIDUnsupported Error = "string as ID is unsupported, please use ExtractMangaID and pass it to methods"
	ErrInvalidID             Error = "invalid ID"
)

func (e Error) Error() string {
	return string(e)
}

func (e ExtractedURL) MustMangaID() any {
	return e["manga"]
}

func (e ExtractedURL) MangaID() (any, error) {
	result, ok := e["manga"]
	if !ok {
		return "", ErrInvalidID
	}

	return result, nil
}

func (e ExtractedURL) ChapterID() (any, error) {
	result, ok := e["chapter"]
	if !ok {
		return "", ErrInvalidID
	}

	return result, nil
}

func (e ExtractedURL) URL() (any, error) {
	result, ok := e["url"]
	if !ok {
		return "", ErrInvalidID
	}

	return result, nil
}
