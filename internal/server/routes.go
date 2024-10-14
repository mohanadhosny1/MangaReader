package server

import (
	"MangaReader/internal/manga/mangafire"
	"MangaReader/pkg/httpClient"
	"github.com/coocood/freecache"
	cache "github.com/gitsight/go-echo-cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func (s *Server) RegisterRoutes() http.Handler {
	c := freecache.NewCache(1024 * 1024)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(cache.New(&cache.Config{}, c))
	e.GET("/", s.HealthHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	group := e.Group("/api")
	group.POST("/search", s.SearchHandler)
	group.POST("/chapter", s.ChapterHandler)
	group.POST("/chapters", s.ChaptersHandler)

	return e
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) SearchHandler(c echo.Context) error {
	searchRequest := new(SearchRequest)
	if err := c.Bind(searchRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "Bad request!"})
	}

	if searchRequest.Query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Query is required!"})
	}

	client, _ := httpClient.NewHttpClient("", time.Second*10, true)
	mf := mangafire.NewMangaFire(client)

	mangas, err := mf.Search(searchRequest.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok", "mangas": mangas})
}

func (s *Server) ChapterHandler(c echo.Context) error {
	dataRequest := new(DataRequest)
	if err := c.Bind(dataRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "Bad request!"})
	}

	client, _ := httpClient.NewHttpClient("", time.Second*10, true)
	mf := mangafire.NewMangaFire(client)

	chapter, err := mf.GetChapter(dataRequest.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok", "chapter": chapter})
}

func (s *Server) ChaptersHandler(c echo.Context) error {
	dataRequest := new(DataRequest)
	if err := c.Bind(dataRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "Bad request!"})
	}

	client, _ := httpClient.NewHttpClient("", time.Second*10, true)
	mf := mangafire.NewMangaFire(client)

	chapters, err := mf.GetManga(dataRequest.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok", "chapters": chapters.Chapters})
}
