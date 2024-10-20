package models

import (
	"database/sql"
	"fmt"
)

type CurrentBookmarkGroup struct {
	BookmarkGroupID int64
}

type CurrentBookmarkGroupRepository struct {
	db *sql.DB
}

func NewCurrentBookmarkGroupRepository(db *sql.DB) *CurrentBookmarkGroupRepository {
	return &CurrentBookmarkGroupRepository{db: db}
}

func (r *CurrentBookmarkGroupRepository) GetCurrentBookmarkGroupId() (int64, error) {
	var bookmarkGroupID int64
	query := "SELECT bookmark_group_id FROM current_bookmark_group WHERE id = 1"
	err := r.db.QueryRow(query).Scan(&bookmarkGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error getting current bookmark group: %w", err)
	}
	return bookmarkGroupID, nil
}

// You might also want to add a method to set the current bookmark group
func (r *CurrentBookmarkGroupRepository) SetCurrentBookmarkGroupId(bookmarkGroupID int64) error {
	query := "UPDATE current_bookmark_group SET bookmark_group_id = ? WHERE id = 1"
	_, err := r.db.Exec(query, bookmarkGroupID)
	if err != nil {
		return fmt.Errorf("error setting current bookmark group: %w", err)
	}
	return nil
}
