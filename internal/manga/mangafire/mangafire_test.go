package mangafire

import (
	"MangaReader/pkg/httpClient"
	"testing"
	"time"
)

func TestMangaFire_Search(t *testing.T) {
	client, _ := httpClient.NewHttpClient("", time.Second*10, true)
	mf := NewMangaFire(client)

	mangas, err := mf.Search("Berserk")
	if err != nil {
		t.Errorf("Failed to search for manga: %s", err)
	}

	if len(mangas) == 0 {
		t.Fatalf("No mangas found")
	}

	t.Logf("Mangas found: %d", len(mangas))
}

func TestMangaFire_GetManga(t *testing.T) {
	client, _ := httpClient.NewHttpClient("", time.Second*10, true)
	mf := NewMangaFire(client)

	info, err := mf.GetManga("m2vv")
	if err != nil {
		t.Fatalf("Failed to search for manga: %s", err)
	}

	if info == nil {
		t.Fatalf("No manga found")
	}

	t.Logf("Manga found: %s", info.Name)
	t.Logf("Poster Found: %s", info.Poster)
	t.Logf("Chapters found: %d", len(info.Chapters))
	t.Logf("Chapters found: %v", info.Chapters[0])
}

func TestMangaFire_GetChapter(t *testing.T) {
	client, _ := httpClient.NewHttpClient("", time.Second*10, true)
	mf := NewMangaFire(client)

	images, err := mf.GetChapter("https://mangafire.to/ajax/read/chapter/3693036")
	if err != nil {
		t.Fatalf("Failed to search for manga: %s", err)
	}

	for i, image := range images {
		t.Logf("Image %d: %s", i+1, image)
	}
}
