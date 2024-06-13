package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/ineersa/blog/models"
	"github.com/ineersa/blog/structs"
	"github.com/ineersa/blog/templates"
	"github.com/ineersa/blog/templates/pages"
	"github.com/ineersa/blog/templates/partials"
)

// indexViewHandler handles a view for the index page.
func (s *Server) indexViewHandler(w http.ResponseWriter, r *http.Request) {
	metadata := structs.Metadata{
		Title:          "Welcome to Ineersa Blog!",
		Keywords:       "blog, ineersa, tech, technical, php, server, python, golang",
		Description:    "Welcome to Ineersa blog!",
		IsNeedToRender: true,
	}

	categories, err := s.categoriesModel.GetCategories()
	if err != nil {
		slog.Error("get categories error", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pageData := models.NewPostsListPageData(r.RequestURI)
	postsList, err := s.postsModel.GetPostsList(pageData)
	if err != nil {
		slog.Error("get posts error", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tagsWithCounts, err := s.tagsModel.GetTagsWithCount(r.RequestURI)
	if err != nil {
		slog.Error("get tags with count error:", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var indexTemplate templ.Component

	if htmx.IsHTMX(r) {
		bodyContent := pages.IndexBodyContent(
			categories,
			postsList,
			pageData,
			tagsWithCounts,
			metadata,
		)
		indexTemplate = bodyContent
	} else {
		bodyContent := pages.IndexBodyContent(categories, postsList, pageData, tagsWithCounts, metadata.SetIsNeedToRender(false))
		indexTemplate = templates.Layout(
			metadata.SetIsNeedToRender(false),
			bodyContent,
		)
	}

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, indexTemplate); err != nil {
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func (s *Server) aboutViewHandler(w http.ResponseWriter, r *http.Request) {
	var aboutTemplate templ.Component
	metadata := structs.Metadata{
		Title:          "About Page!",
		Keywords:       "blog, ineersa, tech, technical, php, server, python, golang, about",
		Description:    "About Page!",
		IsNeedToRender: true,
	}
	if htmx.IsHTMX(r) {
		aboutTemplate = partials.About(metadata)
	} else {
		bodyContent := pages.AboutBodyContent(metadata.SetIsNeedToRender(false))
		aboutTemplate = templates.Layout(
			metadata.SetIsNeedToRender(false),
			bodyContent,
		)
	}

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, aboutTemplate); err != nil {
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send log message.
	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}

func (s *Server) blogPostViewHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	var blogPostTemplate templ.Component
	post, err := s.postsModel.GetPostDetails(slug)

	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("Not found error", "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("get post details error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	metadata := structs.Metadata{
		Title:          post.Title,
		Keywords:       strings.Join(post.Keywords, ","),
		Description:    post.ShortDescription,
		IsNeedToRender: true,
	}

	if htmx.IsHTMX(r) {
		blogPostTemplate = partials.BlogPostBodyContent(post, metadata)
	} else {
		bodyContent := pages.BlogPost(post, metadata.SetIsNeedToRender(false))

		blogPostTemplate = templates.Layout(
			metadata.SetIsNeedToRender(false),
			bodyContent,
		)
	}

	if err := htmx.NewResponse().RenderTempl(r.Context(), w, blogPostTemplate); err != nil {
		slog.Error("render template", "method", r.Method, "status", http.StatusInternalServerError, "path", r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("render page", "method", r.Method, "status", http.StatusOK, "path", r.URL.Path)
}
