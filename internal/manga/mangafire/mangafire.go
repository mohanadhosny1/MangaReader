package mangafire

import (
	"MangaReader/internal/manga"
	"MangaReader/pkg/httpClient"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"strings"
)

type MangaFire struct {
	client *httpClient.HttpClient
}

func NewMangaFire(client *httpClient.HttpClient) manga.Provider {
	return &MangaFire{
		client: client,
	}
}

func (m MangaFire) Search(query string) ([]manga.Search, error) {
	res, err := m.client.Get("https://mangafire.to/ajax/manga/search?keyword="+query, nil)
	if err != nil {
		return nil, err
	}

	body := gjson.Parse(res.Body)
	if res.StatusCode != 200 || body.Get("status").Int() != 200 {
		return nil, manga.ErrStatusNotOK
	}

	htmlContent := body.Get("result.html").String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var mangas []manga.Search
	doc.Find(".unit").Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "")

		var id string
		if strings.Contains(url, ".") {
			id = strings.Split(url, ".")[1]
		}

		mangas = append(mangas, manga.Search{
			Name:   s.Find(".info h6").Text(),
			ID:     id,
			URL:    "https://mangafire.to" + url,
			Poster: s.Find("img").First().AttrOr("src", ""),
		})
	})

	if len(mangas) == 0 {
		return nil, manga.ErrMangaNotFound
	}

	return mangas, nil
}

func (m MangaFire) GetManga(id string) (*manga.Manga, error) {
	res, err := m.client.Get(fmt.Sprintf("https://mangafire.to/ajax/read/%s/chapter/en", id), nil)
	if err != nil {
		return nil, err
	}

	body := gjson.Parse(res.Body)
	if res.StatusCode != 200 || body.Get("status").Int() != 200 {
		return nil, manga.ErrStatusNotOK
	}

	htmlContent := body.Get("result.html").String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var chapters []manga.Chapter
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		info := s.Find("a").First()
		chapters = append(chapters, manga.Chapter{
			Number: info.AttrOr("data-number", ""),
			Name:   info.AttrOr("title", ""),
			ID:     info.AttrOr("data-id", ""),
		})
	})

	return &manga.Manga{Chapters: chapters}, nil
}

func (m MangaFire) GetChapter(id string) ([]string, error) {
	res, err := m.client.Get(fmt.Sprintf("https://mangafire.to/ajax/read/chapter/%s", id), nil)
	if err != nil {
		return nil, err
	}

	body := gjson.Parse(res.Body)
	if res.StatusCode != 200 || body.Get("status").Int() != 200 {
		return nil, manga.ErrStatusNotOK
	}

	var images []string
	data := body.Get("result.images.#.0").Array()
	for _, image := range data {
		images = append(images, image.String())
	}

	return images, nil
}
