package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ineersa/blog/services"
)

const DefaultPageSize = 5
const MinSearchSize = 3
const URIPage = "page"
const URICategory = "category"
const URITags = "tags"
const URISearch = "search"

var uriParams = []string{
	URIPage,
	URICategory,
	URITags,
	URISearch,
}

type PostsModel struct {
	database        *sql.DB
	tagsModel       *TagsModel
	categoriesModel *CategoriesModel
}

func NewPostsModel(database *sql.DB, tagsModel *TagsModel, categoriesModel *CategoriesModel) *PostsModel {
	return &PostsModel{database, tagsModel, categoriesModel}
}

type PostsListItem struct {
	Post Post  `json:"post"`
	Tags []Tag `json:"tags"`
}

type Post struct {
	Title            string    `json:"title"             db:"title"`
	Slug             string    `json:"slug"              db:"slug"`
	Thumbnail        string    `json:"thumbnail"         db:"thumbnail"`
	Color            string    `json:"color"             db:"color"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	CreatedAt        time.Time `json:"created_at"        db:"created_at"`
	PublishedAt      time.Time `json:"published_at"      db:"published_at"`
}

type BlogPostData struct {
	Title            string          `json:"title"             db:"title"`
	Slug             string          `json:"slug"              db:"slug"`
	Thumbnail        string          `json:"thumbnail"         db:"thumbnail"`
	Color            string          `json:"color"             db:"color"`
	CreatedAt        time.Time       `json:"created_at"        db:"created_at"`
	PublishedAt      time.Time       `json:"published_at"      db:"published_at"`
	Keywords         []string        `json:"keywords"          db:"keywords"`
	Content          string          `json:"content"           db:"content"`
	ShortDescription string          `json:"short_description" db:"short_description"`
	Tags             []Tag           `json:"tags"`
	NextArticle      ArticleLinkInfo `json:"next_article"`
	PreviousArticle  ArticleLinkInfo `json:"previous_article"`
}

type ArticleLinkInfo struct {
	IsExist     bool        `json:"is_exist"`
	ArticleInfo ArticleInfo `json:"article_info"`
}

type ArticleInfo struct {
	Title string `json:"title" db:"title"`
	Slug  string `json:"slug"  db:"slug"`
	Color string `json:"color" db:"color"`
}

type PostsListPageData struct {
	Pagination Pagination
	Filters    Filters
	RequestURI string
}

type Pagination struct {
	PagesCount int `json:"pages_count"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
}

type Filters struct {
	Category string `json:"category"`
	Tags     string `json:"tags"`
	Search   string `json:"search"`
}

func NewPostsListPageData(requestURI string) *PostsListPageData {
	pageCount := 0
	currentPage := 1
	parsedURI, err := url.ParseRequestURI(requestURI)
	filters := Filters{
		Category: "",
		Tags:     "",
		Search:   "",
	}
	if err == nil {
		values := parsedURI.Query()
		currentPage, _ = strconv.Atoi(values.Get("page"))
		if currentPage < 1 {
			currentPage = 1
		}
		filters = Filters{
			Category: values.Get("category"),
			Tags:     values.Get("tags"),
			Search:   values.Get("search"),
		}
	}
	pagination := Pagination{
		PagesCount: pageCount,
		Page:       currentPage,
		Limit:      DefaultPageSize,
		Offset:     (currentPage - 1) * DefaultPageSize,
	}

	return &PostsListPageData{
		Pagination: pagination,
		Filters:    filters,
		RequestURI: requestURI,
	}
}

func (pageData *PostsListPageData) GetValueForURI(uriParam string) string {
	switch uriParam {
	case URITags:
		return pageData.Filters.Tags
	case URISearch:
		return pageData.Filters.Search
	case URICategory:
		return pageData.Filters.Category
	case URIPage:
		return strconv.Itoa(pageData.Pagination.Page)
	default:
		return ""
	}
}

func (pageData *PostsListPageData) GetLink(changes map[string]string) string {
	parsedURI, err := url.ParseRequestURI(pageData.RequestURI)
	if err != nil {
		return "/"
	}

	values := parsedURI.Query()
	for _, param := range uriParams {
		values.Del(param)
		paramValue, ok := changes[param]
		if !ok {
			paramValue = pageData.GetValueForURI(param)
		}
		if paramValue != "" {
			values.Add(param, paramValue)
		}
	}

	parsedURI.RawQuery = values.Encode()
	return parsedURI.String()
}

func (pageData *PostsListPageData) formWherePart(tagsDictionary map[int]Tag, categoriesDictionary map[int]Category) (where string, values []any) {
	where = " WHERE p.published = 1 "
	if len(pageData.Filters.Search) >= MinSearchSize {
		searchResponse := services.Search(pageData.Filters.Search)
		if len(searchResponse.IDs) > 0 {
			inClause, args := buildInClause("p.id", searchResponse.IDs)
			where += fmt.Sprintf(" AND %s ", inClause)
			values = append(values, args...)
		} else {
			where += " AND p.title LIKE ? "
			values = append(values, "%"+pageData.Filters.Search+"%")
		}
	}
	if pageData.Filters.Category != "" {
		category := 0
		for categoryID, categoryItem := range categoriesDictionary {
			if pageData.Filters.Category == categoryItem.Slug {
				category = categoryID
			}
		}
		if category > 0 {
			where += " AND p.category_id = ? "
			values = append(values, category)
		}
	}
	if pageData.Filters.Tags != "" {
		tags := strings.Split(pageData.Filters.Tags, ",")
		tagIDs := make([]int, len(tags))
		for tagID, tagItem := range tagsDictionary {
			for key, tag := range tags {
				if tag == tagItem.Slug {
					tagIDs[key] = tagID
				}
			}
		}
		inClause, args := buildInClause("pt.tag_id", tagIDs)
		where += fmt.Sprintf(" AND %s ", inClause)
		values = append(values, args...)
	}

	return where, values
}

func buildInClause(field string, values []int) (inClause string, args []interface{}) {
	placeholders := make([]string, len(values))
	args = make([]interface{}, len(values))
	for i, v := range values {
		placeholders[i] = "?"
		args[i] = v
	}
	inClause = fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", "))
	return inClause, args
}

func (post *Post) GetThumbnailLink() string {
	return "/shared/" + post.Thumbnail
}

func (model *PostsModel) GetCount(wherePart string, whereValues []any) (int, error) {
	count := 0
	selectPart := `
		SELECT
		COUNT(1)
	`
	fromPart := `
		FROM posts AS p
		LEFT JOIN posts_tags AS pt ON pt.post_id = p.id
	`
	groupByPart := ` GROUP BY p.id`

	innerQuery := selectPart + fromPart + wherePart + groupByPart

	query := "SELECT COUNT(1) AS cnt FROM (" + innerQuery + ") AS t"
	row := model.database.QueryRow(query, whereValues...)

	if err := row.Scan(&count); err != nil {
		return count, fmt.Errorf("row scan error %w", err)
	}

	return count, nil
}

func (model *PostsModel) GetPostsList(pageData *PostsListPageData) ([]*PostsListItem, error) {
	var postsList []*PostsListItem
	tagsDictionary, err := model.tagsModel.GetTagsDictionary()
	if err != nil {
		return nil, fmt.Errorf("GetTagsDictionary error: %w", err)
	}
	categoriesDictionary, err := model.categoriesModel.GetCategoriesDictionary()
	if err != nil {
		return nil, fmt.Errorf("GetCategoriesDictionary error: %w", err)
	}

	selectPart := `
		SELECT
		p.title, p.slug, p.thumbnail, p.color,
		p.short_description, p.created_at, p.published_at, p.keywords,
		GROUP_CONCAT(DISTINCT pt_actual.tag_id) AS tag_ids
	`
	fromPart := `
		FROM posts AS p
		LEFT JOIN posts_tags AS pt ON pt.post_id = p.id
		LEFT JOIN posts_tags AS pt_actual ON pt_actual.post_id = p.id
	`
	wherePart, whereValues := pageData.formWherePart(tagsDictionary, categoriesDictionary)

	groupByPart := " GROUP BY p.id "

	orderPart := ` ORDER BY p.published_at DESC, p.created_at DESC `

	limitOffsetPart := fmt.Sprintf(" LIMIT %v OFFSET %v", pageData.Pagination.Limit, pageData.Pagination.Offset)

	query := selectPart + fromPart + wherePart + groupByPart + orderPart + limitOffsetPart
	rows, err := model.database.Query(query, whereValues...)
	if err != nil {
		return nil, fmt.Errorf("GetPostsList error: %w", err)
	}

	count, err := model.GetCount(wherePart, whereValues)
	if err != nil {
		return nil, fmt.Errorf("GetCount error: %w", err)
	}
	pageCount := int(math.Ceil(float64(count) / float64(DefaultPageSize)))
	if pageCount == 0 {
		pageCount = 1
	}
	pageData.Pagination.PagesCount = pageCount

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Info("error closing rows:", "error", err.Error())
		}
	}(rows)

	for rows.Next() {
		var (
			title            string
			slug             string
			thumbnail        string
			color            string
			shortDescription string
			createdAt        time.Time
			publishedAt      time.Time
			keywords         string
			tagIDs           sql.NullString
		)
		if err := rows.Scan(
			&title,
			&slug,
			&thumbnail,
			&color,
			&shortDescription,
			&createdAt,
			&publishedAt,
			&keywords,
			&tagIDs,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		postsList = append(postsList, &PostsListItem{
			Post: Post{
				Title:            title,
				Slug:             slug,
				Thumbnail:        thumbnail,
				Color:            color,
				ShortDescription: shortDescription,
				CreatedAt:        createdAt,
				PublishedAt:      publishedAt,
			},
			Tags: func() []Tag {
				var tags []Tag
				if tagIDs.Valid {
					tagsIDsParsed := strings.Split(tagIDs.String, ",")
					for _, tagID := range tagsIDsParsed {
						tagIDInt, _ := strconv.Atoi(tagID)
						tags = append(tags, tagsDictionary[tagIDInt])
					}
				}

				return tags
			}(),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning rows: %w", err)
	}

	return postsList, nil
}

func (model *PostsModel) GetPostDetails(slugURL string) (BlogPostData, error) {
	var article BlogPostData
	tagsDictionary, err := model.tagsModel.GetTagsDictionary()
	if err != nil {
		return article, fmt.Errorf("GetTagsDictionary error: %w", err)
	}

	query := `
		SELECT
			p.title, p.slug, p.thumbnail, p.color,
			p.created_at, p.published_at, p.keywords, p.content, p.short_description,
			GROUP_CONCAT(pt.tag_id) AS tag_ids
		FROM posts AS p
		LEFT JOIN posts_tags AS pt ON pt.post_id = p.id
		WHERE p.published=1 AND p.slug=?
		GROUP BY p.id`

	row := model.database.QueryRow(query, slugURL)
	var (
		title            string
		slug             string
		thumbnail        string
		color            string
		createdAt        time.Time
		publishedAt      time.Time
		keywords         string
		content          string
		shortDescription string
		tagIDs           sql.NullString
	)
	if err := row.Scan(&title, &slug, &thumbnail, &color, &createdAt, &publishedAt, &keywords, &content, &shortDescription, &tagIDs); err != nil {
		return article, err
	}
	return BlogPostData{
		Title:       title,
		Slug:        slug,
		Thumbnail:   thumbnail,
		Color:       color,
		CreatedAt:   createdAt,
		PublishedAt: publishedAt,
		Keywords: func() []string {
			var arr []string
			_ = json.Unmarshal([]byte(keywords), &arr)
			return arr
		}(),
		Content:          content,
		ShortDescription: shortDescription,
		Tags: func() []Tag {
			var tags []Tag
			if tagIDs.Valid {
				tagsIDsParsed := strings.Split(tagIDs.String, ",")
				for _, tagID := range tagsIDsParsed {
					tagIDInt, _ := strconv.Atoi(tagID)
					tags = append(tags, tagsDictionary[tagIDInt])
				}
			}
			return tags
		}(),
		NextArticle:     model.getNextLinkInfo(slugURL, createdAt, publishedAt),
		PreviousArticle: model.getPreviousLinkInfo(slugURL, createdAt, publishedAt),
	}, nil
}

func (model *PostsModel) getNextLinkInfo(slugURL string, createdAt, publishedAt time.Time) ArticleLinkInfo {
	nextQuery := `
		SELECT p.title, p.slug, p.color
		FROM posts AS p
		WHERE p.published = 1 AND p.slug != ?
		AND p.published_at >= ?
		AND p.created_at >= ?
		ORDER BY p.published_at ASC, p.created_at ASC
		LIMIT 1
	`
	nextRow := model.database.QueryRow(nextQuery, slugURL, publishedAt, createdAt)
	var (
		linkTitle string
		linkSlug  string
		linkColor string
	)
	nextArticleInfo := ArticleLinkInfo{
		IsExist:     false,
		ArticleInfo: ArticleInfo{},
	}
	if err := nextRow.Scan(&linkTitle, &linkSlug, &linkColor); err == nil {
		nextArticleInfo.IsExist = true
		nextArticleInfo.ArticleInfo = ArticleInfo{
			Title: linkTitle,
			Slug:  linkSlug,
			Color: linkColor,
		}
	}

	return nextArticleInfo
}

func (model *PostsModel) getPreviousLinkInfo(slugURL string, createdAt, publishedAt time.Time) ArticleLinkInfo {
	prevQuery := `
		SELECT p.title, p.slug, p.color
		FROM posts AS p
		WHERE p.published = 1 AND p.slug != ?
		AND p.published_at <= ?
		AND p.created_at <= ?
		ORDER BY p.published_at DESC, p.created_at DESC
		LIMIT 1
	`
	prevRow := model.database.QueryRow(prevQuery, slugURL, publishedAt, createdAt)
	var (
		linkTitle string
		linkSlug  string
		linkColor string
	)
	prevArticleInfo := ArticleLinkInfo{
		IsExist:     false,
		ArticleInfo: ArticleInfo{},
	}
	if err := prevRow.Scan(&linkTitle, &linkSlug, &linkColor); err == nil {
		prevArticleInfo.IsExist = true
		prevArticleInfo.ArticleInfo = ArticleInfo{
			Title: linkTitle,
			Slug:  linkSlug,
			Color: linkColor,
		}
	}

	return prevArticleInfo
}
