package models

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type CategoriesModel struct {
	database *sql.DB
}

func NewCategoriesModel(database *sql.DB) *CategoriesModel {
	return &CategoriesModel{database: database}
}

type Category struct {
	ID   int    `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

func (model *CategoriesModel) GetCategoriesDictionary() (map[int]Category, error) {
	categoriesDictionary := make(map[int]Category)
	rows, err := model.database.Query("SELECT id, name, slug FROM categories")
	if err != nil {
		return nil, fmt.Errorf("get categories dict error: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Info("error closing rows: %w", "error", err.Error())
		}
	}(rows)

	for rows.Next() {
		var category Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		categoriesDictionary[category.ID] = category
	}

	return categoriesDictionary, nil
}

func (model *CategoriesModel) GetCategories() ([]*Category, error) {
	var categories []*Category

	rows, err := model.database.Query(
		`SELECT id, name, slug FROM categories`)
	if err != nil {
		return nil, fmt.Errorf("GetTagsWithCount error: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Info("error closing rows", "error", err.Error())
		}
	}(rows)

	for rows.Next() {
		var (
			id   int64
			name string
			slug string
		)
		err := rows.Scan(&id, &name, &slug)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		categories = append(categories, &Category{
			ID:   int(id),
			Name: name,
			Slug: slug,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning rows: %w", err)
	}

	return categories, nil
}
