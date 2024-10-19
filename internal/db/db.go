package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Function to get the database connection
func GetDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Function to initialize the database
func InitDB(path string) error {
	db, err := GetDB(path)

	if err != nil {
		return err
	}

	defer db.Close()

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	// Create bookmark_group table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bookmark_group (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_bookmark_name ON bookmark_group(name);
	`)
	if err != nil {
		return err
	}

	// Create branches table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS branches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bookmark_group_id INTEGER,
			name TEXT NOT NULL,
			branch_alias TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (bookmark_group_id) REFERENCES bookmark_group(id) ON DELETE CASCADE,
			UNIQUE(bookmark_group_id, name)
		);
		CREATE INDEX IF NOT EXISTS idx_branch_name ON branches(name);
		CREATE INDEX IF NOT EXISTS idx_bookmark_group_id ON branches(bookmark_group_id);
	`)
	if err != nil {
		return err
	}

	// Create current bookmark group table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS current_bookmark_group (
        id INTEGER PRIMARY KEY CHECK (id = 1),
        bookmark_group_id INTEGER NOT NULL,
        FOREIGN KEY (bookmark_group_id) REFERENCES bookmark_group(bookmark_group_id) ON DELETE SET NULL
    );
	`)
	if err != nil {
		return err
	}

	return nil
}
