package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Branch struct {
	ID              int64
	BookmarkGroupID int64
	Name            string
	Alias           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (b *Branch) Create(db *sql.DB) error {
	query := `
        INSERT INTO branches (bookmark_group_id, name, branch_alias, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
	now := time.Now()
	result, err := db.Exec(query, b.BookmarkGroupID, b.Name, b.Alias, now, now)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return fmt.Errorf("branch '%s' already exists in this bookmark group", b.Name)
			}
		}
		return fmt.Errorf("error creating branch: %w", err)
	}
	b.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *Branch) ListByBookmarkGroupId(db *sql.DB) ([]Branch, error) {
	query := "SELECT id, name, branch_alias, created_at, updated_at FROM branches WHERE bookmark_group_id = ?"
	rows, err := db.Query(query, b.BookmarkGroupID)
	if err != nil {
		return nil, fmt.Errorf("error querying branches: %w", err)
	}
	defer rows.Close()

	var branches []Branch
	for rows.Next() {
		var branch Branch
		err := rows.Scan(&branch.ID, &branch.Name, &branch.Alias, &branch.CreatedAt, &branch.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		branch.BookmarkGroupID = b.BookmarkGroupID
		branches = append(branches, branch)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return branches, nil
}

func (b *Branch) GetByName(db *sql.DB) error {
	query := "SELECT id, branch_alias, created_at, updated_at FROM branches WHERE bookmark_group_id = ? AND name = ?"
	err := db.QueryRow(query, b.BookmarkGroupID, b.Name).Scan(&b.ID, &b.Alias, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("branch '%s' not found in this bookmark group", b.Name)
		}
		return fmt.Errorf("error getting branch: %w", err)
	}
	return nil
}
