package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/ineersa/blog/models"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "blog",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	var db *sql.DB
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		slog.Error("Failed to open mysql connection!", "details", err.Error())
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	slog.Info("Mysql connected!")

	tagsModel := models.NewTagsModel(db)
	categoriesModel := models.NewCategoriesModel(db)
	postsModel := models.NewPostsModel(db, tagsModel, categoriesModel)

	server := NewServer(
		tagsModel,
		categoriesModel,
		postsModel,
	)

	// Run your server.
	if err := server.Run(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}

}
