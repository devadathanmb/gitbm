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

type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(b *Branch) error {
	query := `
        INSERT INTO branches (bookmark_group_id, name, branch_alias, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
	now := time.Now()
	result, err := r.db.Exec(query, b.BookmarkGroupID, b.Name, b.Alias, now, now)
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

func (r *BranchRepository) ListByBookmarkGroupId(bookmarkGroupID int64) ([]Branch, error) {
	query := "SELECT id, name, branch_alias, created_at, updated_at FROM branches WHERE bookmark_group_id = ?"
	rows, err := r.db.Query(query, bookmarkGroupID)
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
		branch.BookmarkGroupID = bookmarkGroupID
		branches = append(branches, branch)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return branches, nil
}

func (r *BranchRepository) GetByName(bookmarkGroupID int64, name string) (*Branch, error) {
	query := "SELECT id, branch_alias, created_at, updated_at FROM branches WHERE bookmark_group_id = ? AND name = ?"
	b := &Branch{BookmarkGroupID: bookmarkGroupID, Name: name}
	err := r.db.QueryRow(query, bookmarkGroupID, name).Scan(&b.ID, &b.Alias, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("branch '%s' not found in this bookmark group", name)
		}
		return nil, fmt.Errorf("error getting branch: %w", err)
	}
	return b, nil
}

func (r *BranchRepository) Remove(bookmarkGroupID int64, name string) error {
	query := "DELETE FROM branches WHERE bookmark_group_id = ? AND name = ?"
	_, err := r.db.Exec(query, bookmarkGroupID, name)
	if err != nil {
		return fmt.Errorf("error removing branch: %w", err)
	}
	return nil
}
