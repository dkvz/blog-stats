package db

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type DbSqlite struct {
	db  *sql.DB
	ctx context.Context
}

func NewDBSqlite(path string) (*DbSqlite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return &DbSqlite{
		db:  db,
		ctx: context.Background(),
	}, nil
}

func (dbs *DbSqlite) AllPublishedArticleIds() ([]uint, error) {
	rows, err := dbs.db.QueryContext(
		dbs.ctx,
		`SELECT id FROM articles WHERE published=1 ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// If this seems overcomplicated, that's because it is.
	// People often use the sqlx package to add annotations
	// to entity structs making the conversion automatic.
	// Otherwise we have to use Scan.
	var ids []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (dbs *DbSqlite) ArticleContentById(id uint) (*string, error) {
	row := dbs.db.QueryRowContext(
		dbs.ctx,
		`SELECT content FROM articles WHERE id=?`, id,
	)

	var content *string
	err := row.Scan(&content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
