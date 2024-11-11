package manga

const (
	ErrUnimplemented         Error = "unimplemented method"
	ErrForbidden             Error = "forbidden"
	ErrMangaNotFound         Error = "manga not found"
	ErrChapterNotFound       Error = "chapter not found"
	ErrPaidChapter           Error = "chapter is paid"
	ErrAuthorizationRequired Error = "authorization required"
)

func (e Error) Error() string {
	return string(e)
}
