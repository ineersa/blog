package main

import (
	"fmt"
	"github.com/ineersa/blog/models"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	gowebly "github.com/gowebly/helpers"
)

const ReadTimeout = 5
const WriteTimeout = 10

type Server struct {
	tagsModel       *models.TagsModel
	categoriesModel *models.CategoriesModel
	postsModel      *models.PostsModel
}

func NewServer(tagsModel *models.TagsModel, categoriesModel *models.CategoriesModel, postsModel *models.PostsModel) *Server {
	return &Server{
		tagsModel:       tagsModel,
		categoriesModel: categoriesModel,
		postsModel:      postsModel,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.indexViewHandler)
	mux.HandleFunc("GET /about", s.aboutViewHandler)
	mux.HandleFunc("GET /post/{slug}", s.blogPostViewHandler)

	return mux
}

func (s *Server) Run() error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "7001"))
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      s.Handler(),
		ReadTimeout:  ReadTimeout * time.Second,
		WriteTimeout: WriteTimeout * time.Second,
	}

	slog.Info("Starting server...", "port", port)

	err = httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
