package models

import (
	"database/sql"
	"fmt"
)

type BookmarkGroup struct {
	ID   int64
	Name string
}

// Create a new bookmark group and set it as the current bookmark group
func (bg *BookmarkGroup) Create(db *sql.DB) error {
	// Start a transaction
	tx, err := db.Begin()
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
func (bg *BookmarkGroup) List(db *sql.DB) ([]BookmarkGroup, error) {
	query := "SELECT id, name FROM bookmark_group"
	rows, err := db.Query(query)
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
func (bg *BookmarkGroup) GetCurrent(db *sql.DB) error {
	query := `
		SELECT bg.id, bg.name 
		FROM bookmark_group bg
		JOIN current_bookmark_group cbg ON bg.id = cbg.bookmark_group_id
		WHERE cbg.id = 1
	`
	err := db.QueryRow(query).Scan(&bg.ID, &bg.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no current bookmark group set")
		}
		return fmt.Errorf("error getting current bookmark group: %w", err)
	}
	return nil
}
