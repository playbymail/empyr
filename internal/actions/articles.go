// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"database/sql"
	"github.com/playbymail/empyr/internal/models"
	"time"
)

type ArticlesFacade struct {
	db *sql.DB
}

func (f *ArticlesFacade) GetPublishedArticles() ([]models.Article, error) {

	sqlStmt := `
SELECT id, title, slug, published, date_published, date_updated
FROM articles
WHERE published = 1
ORDER BY id DESC;
`

	rows, err := f.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var article models.Article

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Published,
			&article.DatePublished,
			&article.DateUpdated,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)

	}

	return articles, nil
}

func (f *ArticlesFacade) CreateArticle(article *models.Article) error {
	sqlStmt := `
        INSERT INTO articles (title, slug, published, date_published, date_updated)
        VALUES (?, ?, ?, ?, ?);
    `
	result, err := f.db.Exec(sqlStmt,
		article.Title,
		article.Slug,
		article.Published,
		article.DatePublished,
		article.DateUpdated)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	article.ID = int(id)
	return nil
}

func (f *ArticlesFacade) GetArticleByID(id int) (*models.Article, error) {
	sqlStmt := `
        SELECT id, title, slug, published, date_published, date_updated
        FROM articles
        WHERE id = ?;
    `
	var article models.Article
	err := f.db.QueryRow(sqlStmt, id).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Published,
		&article.DatePublished,
		&article.DateUpdated)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (f *ArticlesFacade) GetArticleBySlug(slug string) (*models.Article, error) {
	sqlStmt := `
        SELECT id, title, slug, published, date_published, date_updated
        FROM articles
        WHERE slug = ?;
    `
	var article models.Article
	err := f.db.QueryRow(sqlStmt, slug).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Published,
		&article.DatePublished,
		&article.DateUpdated)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (f *ArticlesFacade) UpdateArticle(article *models.Article) error {
	sqlStmt := `
        UPDATE articles 
        SET title = ?, slug = ?, published = ?, date_updated = ?
        WHERE id = ?;
    `
	_, err := f.db.Exec(sqlStmt,
		article.Title,
		article.Slug,
		article.Published,
		time.Now(),
		article.ID)
	return err
}

func (f *ArticlesFacade) PublishArticle(id int) error {
	sqlStmt := `
        UPDATE articles 
        SET published = 1, date_published = ?, date_updated = ?
        WHERE id = ?;
    `
	now := time.Now()
	_, err := f.db.Exec(sqlStmt, now, now, id)
	return err
}

func (f *ArticlesFacade) UnpublishArticle(id int) error {
	sqlStmt := `
        UPDATE articles 
        SET published = 0, date_updated = ?
        WHERE id = ?;
    `
	_, err := f.db.Exec(sqlStmt, time.Now(), id)
	return err
}

func (f *ArticlesFacade) DeleteArticle(id int) error {
	sqlStmt := `DELETE FROM articles WHERE id = ?;`
	_, err := f.db.Exec(sqlStmt, id)
	return err
}

func (f *ArticlesFacade) GetDraftArticles() ([]models.Article, error) {
	sqlStmt := `
        SELECT id, title, slug, published, date_published, date_updated
        FROM articles
        WHERE published = 0
        ORDER BY id DESC;
    `
	rows, err := f.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Published,
			&article.DatePublished,
			&article.DateUpdated)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (f *ArticlesFacade) GetRecentArticles(limit int) ([]models.Article, error) {
	sqlStmt := `
        SELECT id, title, slug, published, date_published, date_updated
        FROM articles
        WHERE published = 1
        ORDER BY date_published DESC
        LIMIT ?;
    `
	rows, err := f.db.Query(sqlStmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Published,
			&article.DatePublished,
			&article.DateUpdated)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
