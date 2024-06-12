package models

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

type TagsModel struct {
	database *sql.DB
}

func NewTagsModel(database *sql.DB) *TagsModel {
	return &TagsModel{database: database}
}

type TagWithCount struct {
	Tag      Tag  `json:"tag"`
	Count    int  `json:"count" db:"count"`
	Selected bool `json:"selected"`
}

type Tag struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

func (model *TagsModel) GetTagsDictionary() (map[int]Tag, error) {
	tagsDictionary := make(map[int]Tag)
	rows, err := model.database.Query("SELECT id, name, slug FROM tags")
	if err != nil {
		return nil, fmt.Errorf("GetTagsDictionary error: %v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("error closing rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Slug,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		tagsDictionary[int(tag.ID)] = tag
	}

	return tagsDictionary, nil
}

// GetTagsWithCount retrieves tags and count of posts for them, as a return it gives back json.Marshal(tagsWithCount)
func (model *TagsModel) GetTagsWithCount(RequestURI string) ([]*TagWithCount, error) {
	var tagsWithCount []*TagWithCount
	selectedTags := make(map[string]bool)
	parsedURI, err := url.ParseRequestURI(RequestURI)
	if err == nil {
		values := parsedURI.Query()
		tags := values.Get("tags")
		if len(tags) > 0 {
			for _, slug := range strings.Split(tags, ",") {
				selectedTags[slug] = true
			}
		}
	}

	rows, err := model.database.Query(
		`
			SELECT t.id, t.name, t.slug, tags_counts.cnt AS count
			FROM tags AS t
			LEFT JOIN (
				SELECT pt.tag_id, COUNT(1) as cnt
				FROM posts_tags AS pt
				GROUP BY pt.tag_id
			) AS tags_counts ON t.id=tags_counts.tag_id
		`)
	if err != nil {
		return nil, fmt.Errorf("GetTagsWithCount error: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("error closing rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var (
			id   int64
			name string
			slug string
			cnt  sql.NullInt32
		)
		if err := rows.Scan(
			&id,
			&name,
			&slug,
			&cnt,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		tagsWithCount = append(tagsWithCount, &TagWithCount{
			Tag: Tag{ID: id, Name: name, Slug: slug},
			Count: func() int { // or we can do IFNULL() in MySQL
				if cnt.Valid {
					return int(cnt.Int32)
				} else {
					return 0
				}
			}(),
			Selected: func() bool {
				_, ok := selectedTags[slug]
				if ok {
					return true
				}
				return false
			}(),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning rows: %v", err)
	}

	return tagsWithCount, nil
}
