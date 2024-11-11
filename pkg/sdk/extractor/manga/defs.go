package manga

const (
	ErrUnimplemented         Error = "unimplemented method"
	ErrForbidden             Error = "forbidden"
	ErrMangaNotFound         Error = "manga not found"
	ErrChapterNotFound       Error = "chapter not found"
	ErrPaidChapter           Error = "chapter is paid"
	ErrAuthorizationRequired Error = "authorization required"
	ErrInvalidURLFormat      Error = "invalid URL format"
	ErrURLeqID               Error = "URL is ID"
)

func (e Error) Error() string {
	return string(e)
}
