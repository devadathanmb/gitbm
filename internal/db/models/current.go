package models

import (
	"database/sql"
	"fmt"
)

func GetCurrentBookmarkGroupId(db *sql.DB) (int64, error) {
	query := "SELECT bookmark_group_id FROM current_bookmark_group WHERE id = 1"

	var currentBookmarkGroupId int64
	err := db.QueryRow(query).Scan(&currentBookmarkGroupId)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error getting current bookmark group: %w", err)
	}

	return currentBookmarkGroupId, nil
}
