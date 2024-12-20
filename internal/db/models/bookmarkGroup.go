package models

import (
	"database/sql"
	"fmt"
)

type BookmarkGroup struct {
	ID   int64
	Name string
}

type BookmarkGroupRepository struct {
	db *sql.DB
}

func NewBookmarkGroupRepository(db *sql.DB) *BookmarkGroupRepository {
	return &BookmarkGroupRepository{db: db}
}

// Create a new bookmark group and set it as the current bookmark group
func (r *BookmarkGroupRepository) Create(bg *BookmarkGroup) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert the new bookmark group
	result, err := tx.Exec("INSERT INTO bookmark_group (name) VALUES (?)", bg.Name)
	if err != nil {
		return fmt.Errorf("error inserting bookmark group: %w", err)
	}

	// Get the ID of the newly inserted bookmark group
	bg.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	// Check if current bookmark group exists, if not create it
	var currentBookmarkGroupId int64
	err = tx.QueryRow("SELECT bookmark_group_id FROM current_bookmark_group WHERE id = 1").Scan(&currentBookmarkGroupId)
	if err == sql.ErrNoRows {
		// Insert new current bookmark group
		_, err = tx.Exec("INSERT INTO current_bookmark_group (id, bookmark_group_id) VALUES (1, ?)", bg.ID)
	} else if err == nil {
		// Update existing current bookmark group
		_, err = tx.Exec("UPDATE current_bookmark_group SET bookmark_group_id = ? WHERE id = 1", bg.ID)
	}
	if err != nil {
		return fmt.Errorf("error updating current bookmark group: %w", err)
	}

	return tx.Commit()
}

// List all bookmark groups
func (r *BookmarkGroupRepository) List() ([]BookmarkGroup, error) {
	query := "SELECT id, name FROM bookmark_group"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying bookmark groups: %w", err)
	}
	defer rows.Close()

	var bookmarkGroups []BookmarkGroup
	for rows.Next() {
		var group BookmarkGroup
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		bookmarkGroups = append(bookmarkGroups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return bookmarkGroups, nil
}

// Get current bookmark group
func (r *BookmarkGroupRepository) GetCurrent() (*BookmarkGroup, error) {
	query := `
		SELECT bg.id, bg.name 
		FROM bookmark_group bg
		JOIN current_bookmark_group cbg ON bg.id = cbg.bookmark_group_id
		WHERE cbg.id = 1
	`
	var bg BookmarkGroup
	err := r.db.QueryRow(query).Scan(&bg.ID, &bg.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no current bookmark group set")
		}
		return nil, fmt.Errorf("error getting current bookmark group: %w", err)
	}
	return &bg, nil
}

func (r *BookmarkGroupRepository) Delete(bookmarkGroupName string) error {
	query := "DELETE FROM bookmark_group WHERE name = ?"
	_, err := r.db.Exec(query, bookmarkGroupName)
	return err
}

func (r *BookmarkGroupRepository) GetByName(name string) (*BookmarkGroup, error) {
	query := "SELECT id, name FROM bookmark_group WHERE name = ?"
	var bg BookmarkGroup
	err := r.db.QueryRow(query, name).Scan(&bg.ID, &bg.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bookmark group '%s' not found", name)
		}
		return nil, fmt.Errorf("error getting bookmark group: %w", err)
	}
	return &bg, nil
}
