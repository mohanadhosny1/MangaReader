package manga

import "fmt"

var (
	ErrMangaNotFound = fmt.Errorf("manga not found")
	ErrStatusNotOK   = fmt.Errorf("status code not OK")
)
