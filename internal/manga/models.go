package manga

type Provider interface {
	Search(query string) ([]Search, error)
	GetManga(url string) (*Manga, error)
	GetChapter(url string) ([]string, error)
}

type Search struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	URL    string `json:"url"`
	Poster string `json:"poster"`
}

type Chapter struct {
	Number string `json:"number"`
	Name   string `json:"name"`
	ID     string `json:"id"`

	Date string `json:"date,omitempty"`
}

type Manga struct {
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	Poster   string    `json:"poster"`
	Chapters []Chapter `json:"chapters"`
}
